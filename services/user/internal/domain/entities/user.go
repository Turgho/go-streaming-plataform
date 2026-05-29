package entities

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID           string    `bson:"_id"`
	Username     string    `bson:"username"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
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
		return nil, fmt.Errorf("user inválido: %s", strings.Join(errs, ", "))
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
