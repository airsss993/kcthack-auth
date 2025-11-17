package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"

	"github.com/kcthack-auth/internal/domain"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (t *SessionRepo) SaveSession(ctx context.Context, session *domain.Session) error {
	tokenHash := sha256.Sum256([]byte(session.Token))

	query := `INSERT INTO users_sessions (id, user_id, token_hash,expires_at) VALUES ($1, $2, $3, $4)`

	_, err := t.db.ExecContext(ctx, query, session.ID, session.UserID, tokenHash, session.ExpiresAt)
	return err
}

func (t *SessionRepo) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	var session domain.Session

	tokenHash := sha256.Sum256([]byte(token))

	query := `SELECT id, user_id, token_hash, expires_at, created_at FROM users_sessions WHERE token_hash=$1`

	err := t.db.QueryRowContext(ctx, query, tokenHash).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &session, nil
}

func (t *SessionRepo) DeleteByToken(ctx context.Context, token string) error {
	tokenHash := sha256.Sum256([]byte(token))

	query := `DELETE FROM users_sessions WHERE token_hash=$1`

	_, err := t.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (t *SessionRepo) DeleteAllByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM users_sessions WHERE user_id=$1`

	_, err := t.db.ExecContext(ctx, query, userID)
	return err
}
