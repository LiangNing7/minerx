package user

//go:generate mockgen -destination mock_user.go -package user minerx/internal/usercenter/biz/v1/user UserBiz

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/LiangNing7/goutils/pkg/authn"
	"github.com/LiangNing7/goutils/pkg/core"
	"github.com/LiangNing7/goutils/pkg/log"
	"github.com/LiangNing7/goutils/pkg/store/where"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/LiangNing7/minerx/internal/pkg/contextx"
	"github.com/LiangNing7/minerx/internal/pkg/known"
	validationutil "github.com/LiangNing7/minerx/internal/pkg/util/validation"
	"github.com/LiangNing7/minerx/internal/usercenter/model"
	"github.com/LiangNing7/minerx/internal/usercenter/pkg/conversion"
	"github.com/LiangNing7/minerx/internal/usercenter/store"
	v1 "github.com/LiangNing7/minerx/pkg/api/usercenter/v1"
)

// UserBiz defines the interface that contains methods for handling user requests.
type UserBiz interface {
	// Create creates a new user based on the provided request parameters.
	Create(ctx context.Context, rq *v1.CreateUserRequest) (*v1.CreateUserResponse, error)

	// Update updates an existing user based on the provided request parameters.
	Update(ctx context.Context, rq *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error)

	// Delete removes one or more users based on the provided request parameters.
	Delete(ctx context.Context, rq *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error)

	// Get retrieves the details of a specific user based on the provided request parameters.
	Get(ctx context.Context, rq *v1.GetUserRequest) (*v1.GetUserResponse, error)

	// List retrieves a list of users and their total count based on the provided request parameters.
	List(ctx context.Context, rq *v1.ListUserRequest) (*v1.ListUserResponse, error)

	// UserExpansion defines additional methods for extended user operations, if needed.
	UserExpansion
}

// UserExpansion defines additional methods for user operations.
type UserExpansion interface {
	// UpdatePassword updates the password for a user based on the provided request.
	UpdatePassword(ctx context.Context, rq *v1.UpdatePasswordRequest) (*v1.UpdatePasswordResponse, error)
}

// userBiz is the implementation of the UserBiz.
type userBiz struct {
	store store.IStore
}

// Ensure that *userBiz implements the UserBiz.
var _ UserBiz = (*userBiz)(nil)

// New creates and returns a new instance of *userBiz.
func New(store store.IStore) *userBiz {
	return &userBiz{store: store}
}

// Create implements the Create method of the UserBiz.
func (b *userBiz) Create(ctx context.Context, rq *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	var userM model.UserM

	// Copy request data to the User model.
	_ = core.Copy(&userM, rq)

	// Start a transaction for creating the user and secret.
	err := b.store.TX(ctx, func(ctx context.Context) error {
		// Attempt to create the user in the data store.
		if err := b.store.User().Create(ctx, &userM); err != nil {
			// Handle duplicate entry error for username.
			match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error())
			if match {
				return v1.ErrorUserAlreadyExists("user %q already exists", userM.Username)
			}
			return v1.ErrorUserCreateFailed("create user failed: %s", err.Error())
		}

		// Create a secret for the newly created user.
		secretM := &model.SecretM{
			UserID:      userM.UserID,
			Name:        "generated",
			Expires:     0,
			Description: "automatically generated when user is created",
		}
		if err := b.store.Secret().Create(ctx, secretM); err != nil {
			return v1.ErrorSecretCreateFailed("create secret failed: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return nil, err // Return any error from the transaction.
	}

	return &v1.CreateUserResponse{UserID: userM.UserID}, nil
}

// Update implements the Update method of the UserBiz.
func (b *userBiz) Update(ctx context.Context, rq *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	// Update fields if provided in the request.
	if rq.Nickname != nil {
		userM.Nickname = rq.GetNickname()
	}
	if rq.Email != nil {
		userM.Email = rq.GetEmail()
	}
	if rq.Phone != nil {
		userM.Phone = rq.GetPhone()
	}

	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &v1.UpdateUserResponse{}, nil
}

// Delete implements the Delete method of the UserBiz.
func (b *userBiz) Delete(ctx context.Context, rq *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	userID := contextx.UserID(ctx)
	// Limit access to authorized users only.
	if validationutil.IsAdminUser(contextx.UserID(ctx)) {
		userID = rq.GetUserID()
	}

	if err := b.store.User().Delete(ctx, where.F("userID", userID)); err != nil {
		return nil, err
	}

	return &v1.DeleteUserResponse{}, nil
}

// Get implements the Get method of the UserBiz.
func (b *userBiz) Get(ctx context.Context, rq *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	userID := contextx.UserID(ctx)
	// Limit access to authorized users only.
	if validationutil.IsAdminUser(contextx.UserID(ctx)) {
		userID = rq.GetUserID()
	}

	userM, err := b.store.User().Get(ctx, where.F("userID", userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrorUserNotFound(err.Error()) // Return an error if the user is not found.
		}
		return nil, err // Return any other error encountered.
	}

	return &v1.GetUserResponse{User: conversion.UserMToUserV1(userM)}, nil
}

// List implements the List method of the UserBiz.
func (b *userBiz) List(ctx context.Context, rq *v1.ListUserRequest) (*v1.ListUserResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// Set the maximum concurrency limit using the constant MaxConcurrency
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// Use goroutines to improve API performance
	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				converted := conversion.UserMToUserV1(user)

				// Retrieve the count of secrets for each user.
				count, _, err := b.store.Secret().List(ctx, where.F("userID", user.UserID))
				if err != nil {
					log.W(ctx).Errorw(err, "Failed to list secrets")
					return err // Return any error encountered.
				}
				converted.Secrets = count

				m.Store(user.ID, converted)

				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw(err, "Failed to wait all function calls returned")
		return nil, err
	}

	users := make([]*v1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*v1.User))
	}

	return &v1.ListUserResponse{Total: count, Users: users}, nil
}

// UpdatePassword updates the password for a user based on the provided request.
func (b *userBiz) UpdatePassword(ctx context.Context, rq *v1.UpdatePasswordRequest) (*v1.UpdatePasswordResponse, error) {
	// Retrieve the user by username.
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err // Return any error encountered.
	}

	// Compare the old password with the stored password.
	if err := authn.Compare(userM.Password, rq.GetOldPassword()); err != nil {
		return nil, v1.ErrorUserLoginFailed("password incorrect") // Return an error if the old password is incorrect.
	}
	// Encrypt the new password.
	userM.Password, _ = authn.Encrypt(rq.GetNewPassword())

	return &v1.UpdatePasswordResponse{}, b.store.User().Update(ctx, userM) // Update the user's password in the data store.
}
