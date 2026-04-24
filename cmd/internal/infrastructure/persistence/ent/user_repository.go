package persistence

import (
	"context"
	"fmt"

	domain "github.com/dannegm/anubix-server/cmd/internal/domain/user"
	"github.com/dannegm/anubix-server/ent"
	"github.com/dannegm/anubix-server/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func toDomain(e *ent.User) *domain.User {
	return &domain.User{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	users, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = toDomain(u)
	}
	return result, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomain(u), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomain(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	created, err := r.client.User.Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomain(created), nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	updated, err := r.client.User.UpdateOneID(u.ID).
		SetName(u.Name).
		SetEmail(u.Email).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}
	return toDomain(updated), nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
