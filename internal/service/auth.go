package service

import (
	"context"
	"todo-demo/api/v1/dto/auth"
)

type AuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) error
	Login(ctx context.Context, req *auth.LoginRequest) (string, error)
}
