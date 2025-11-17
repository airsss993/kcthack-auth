package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager interface {
	NewAccess(userID, role string, ttl time.Duration) (string, error)
	NewRefresh() string
	Validate(tokenString string) (*TokenClaims, error)
}

type Manager struct {
	signingKey string
}

type TokenClaims struct {
	UserID string
	Role   string
	ExpAt  float64
}

func NewManager(signingKey string) *Manager {
	return &Manager{signingKey: signingKey}
}

func (m *Manager) NewAccess(userID, role string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(ttl).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) NewRefresh() string {
	token := uuid.NewString()
	return token
}

func (m *Manager) Validate(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(m.signingKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to get claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid userID claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role claim")
	}

	expAt, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp claim")
	}

	tokenClaims := TokenClaims{
		UserID: userID,
		Role:   role,
		ExpAt:  expAt,
	}
	return &tokenClaims, nil
}
