package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kcthack-auth/internal/domain"
)

type AuthPSQL struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthPSQL {
	return &AuthPSQL{db: db}
}

func (a *AuthPSQL) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, first_name, last_name, role, email, tg_name, birth_date, bio, pass_hash, is_verified, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := a.db.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Role, user.Email, user.TgName, user.BirthDate, user.BIO, user.PassHash, user.IsVerified, user.UpdatedAt)
	return err
}

func (a *AuthPSQL) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, first_name, last_name, role, email, tg_name, birth_date, bio, pass_hash FROM users WHERE email=$1`

	err := a.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email, &user.TgName, &user.BirthDate, &user.BIO, &user.PassHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AuthPSQL) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, first_name, last_name, role, email, tg_name, birth_date, bio, pass_hash FROM users WHERE id=$1`

	err := a.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role, &user.Email, &user.TgName, &user.BirthDate, &user.BIO, &user.PassHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *AuthPSQL) UpdatePassword(ctx context.Context, userID string, passHash string) error {
	query := `UPDATE users SET pass_hash=$1 WHERE id=$2`

	_, err := a.db.ExecContext(ctx, query, passHash, userID)
	return err
}

func (a *AuthPSQL) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET first_name=$1, last_name=$2, email=$3, tg_name=$4, birth_date=$5, bio=$6, is_verified=$7, updated_at=$8 WHERE id=$9`

	_, err := a.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.TgName, user.BirthDate, user.BIO, user.IsVerified, user.UpdatedAt, user.ID)
	return err
}

func (a *AuthPSQL) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := a.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
