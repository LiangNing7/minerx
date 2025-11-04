package usercenter

import v1 "github.com/LiangNing7/minerx/pkg/api/usercenter/v1"

type GetSecretRequest struct {
	UserID string
	Name   string
}

type GetSecretResponse = v1.Secret
