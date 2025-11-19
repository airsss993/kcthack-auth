package domain

import "time"

const (
	Participant = "participant"
	Partner     = "partner"
	Admin       = "admin"
)

type User struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	Role       string
	TgName     string
	BirthDate  time.Time
	BIO        string
	PassHash   string
	IsVerified bool
	UpdatedAt  time.Time
	CreatedAt  time.Time
}
