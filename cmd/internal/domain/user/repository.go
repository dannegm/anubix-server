package user

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	Delete(ctx context.Context, id int) error
}
