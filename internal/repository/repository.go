package repository

import (
	"context"

	"github.com/kcthack-auth/internal/domain"
)

type AuthRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdatePassword(ctx context.Context, userID string, passHash string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type SessionRepository interface {
	SaveSession(ctx context.Context, session *domain.Session) error
	FindByToken(ctx context.Context, token string) (*domain.Session, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
