package auth

import "context"

type Repository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	AuthHash string `json:"-"`
	Salt     string `json:"-"`
}
