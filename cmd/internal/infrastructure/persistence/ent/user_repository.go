package persistence

import (
	"context"

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

func toDomainUser(e *ent.User) *domain.User {
	return &domain.User{
		ID:       e.ID,
		Email:    e.Email,
		AuthHash: e.AuthHash,
		Salt:     e.Salt,
	}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	users, err := r.client.User.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = toDomainUser(u)
	}
	return result, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u, err := r.client.User.Query().Where(user.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	created, err := r.client.User.Create().
		SetEmail(u.Email).
		SetAuthHash(u.AuthHash).
		SetSalt(u.Salt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(created), nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	updated, err := r.client.User.UpdateOneID(u.ID).
		SetEmail(u.Email).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(updated), nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
