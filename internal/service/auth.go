package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kcthack-auth/internal/domain"
	"github.com/kcthack-auth/internal/repository"
	"github.com/kcthack-auth/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repository.AuthRepository
	srepo      repository.SessionRepository
	tm         auth.JWTManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo repository.AuthRepository, srepo repository.SessionRepository, tm auth.JWTManager, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		srepo:      srepo,
		tm:         tm,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (a *AuthService) Register(ctx context.Context, req RegisterReq) (*AuthResp, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email field cannot be empty")
	}

	if req.FirstName == "" {
		return nil, fmt.Errorf("first name field cannot be empty")
	}

	if req.LastName == "" {
		return nil, fmt.Errorf("last name field cannot be empty")
	}

	exists, err := a.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email: %s already exists", req.Email)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password – user: %v, err: %w", req.Email, err)
	}

	user := domain.User{
		ID:         uuid.NewString(),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Role:       domain.Participant,
		PassHash:   string(passHash),
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	if err := a.repo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("failed to register – user: %v, err: %w", req.Email, err)
	}

	accessToken, err := a.tm.NewAccess(user.ID, user.Role, a.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := a.tm.NewRefresh()

	authResp := AuthResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(a.accessTTL),
	}

	if err := a.srepo.SaveSession(ctx, &domain.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.refreshTTL),
	}); err != nil {
		return nil, fmt.Errorf("failed to save user session: %w", err)
	}

	return &authResp, nil
}

func (a *AuthService) Login(ctx context.Context, req *LoginReq) (*AuthResp, error) {
	if req.Email == "" {
		return nil, fmt.Errorf("email field cannot be empty")
	}

	if req.Password == "" {
		return nil, fmt.Errorf("password field cannot be empty")
	}

	user, err := a.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with email: %s, err: %w", req.Email, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	accessToken, err := a.tm.NewAccess(user.ID, user.Role, a.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := a.tm.NewRefresh()

	authResp := AuthResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(a.accessTTL),
	}

	if err := a.srepo.SaveSession(ctx, &domain.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.refreshTTL),
	}); err != nil {
		return nil, fmt.Errorf("failed to save user session: %w", err)
	}

	return &authResp, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, token string) (*AuthResp, error) {
	session, err := a.srepo.FindByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to find user session: %w", err)
	}

	user, err := a.repo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	accessToken, err := a.tm.NewAccess(user.ID, user.Role, a.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := a.tm.NewRefresh()

	authResp := AuthResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(a.accessTTL),
	}

	if err := a.srepo.SaveSession(ctx, &domain.Session{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.refreshTTL),
	}); err != nil {
		return nil, fmt.Errorf("failed to save user session: %w", err)
	}

	return &authResp, nil
}

func (a *AuthService) Logout(ctx context.Context, token string) error {
	return a.srepo.DeleteByToken(ctx, token)
}
