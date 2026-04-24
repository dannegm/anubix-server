package vault

import "context"

type Repository interface {
	FindAll(ctx context.Context, userID string) ([]*Vault, error)
	FindByID(ctx context.Context, id string) (*Vault, error)
	Create(ctx context.Context, v *Vault) (*Vault, error)
	Update(ctx context.Context, v *Vault) (*Vault, error)
	Delete(ctx context.Context, id string) error
}
