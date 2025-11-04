package usercenter

import (
	usercenterv1 "github.com/LiangNing7/minerx/pkg/api/usercenter/v1"
)

type GetSecretRequest struct {
	Username string
	Name     string
}

type GetSecretResponse usercenterv1.Secret
