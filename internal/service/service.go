package service

import (
	"context"
	"time"
)

type Services struct {
	AuthService AuthService
}

func NewServices(auth AuthService) *Services {
	return &Services{AuthService: auth}
}

type RegisterReq struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type LoginReq struct {
	Email    string
	Password string
}

type AuthResp struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type Auth interface {
	Register(ctx context.Context, req *RegisterReq) (*AuthResp, error)
	Login(ctx context.Context, req *LoginReq) (*AuthResp, error)
	RefreshToken(ctx context.Context, token string) (*AuthResp, error)
	Logout(ctx context.Context, token string) error
}
