package entities

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id string, username, email, passwordHash string) (*User, error) {
	var errs []string

	if id == "" {
		errs = append(errs, "UUID está vazio")
	}
	if username == "" {
		errs = append(errs, "USERNAME é obrigatório")
	}
	if email == "" || !strings.Contains(email, "@") {
		errs = append(errs, "EMAIL é obrigatório e deve ser válido")
	}
	if passwordHash == "" {
		errs = append(errs, "SENHA é obrigatório")
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	now := time.Now().UTC()

	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
