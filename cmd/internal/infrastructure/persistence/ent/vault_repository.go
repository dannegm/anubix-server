package persistence

import (
	"context"

	domain "github.com/dannegm/anubix-server/cmd/internal/domain/vault"
	"github.com/dannegm/anubix-server/ent"
	"github.com/dannegm/anubix-server/ent/user"
	"github.com/dannegm/anubix-server/ent/vault"
)

type VaultRepository struct {
	client *ent.Client
}

func NewVaultRepository(client *ent.Client) *VaultRepository {
	return &VaultRepository{client: client}
}

func toDomainVault(e *ent.Vault) *domain.Vault {
	return &domain.Vault{
		ID:                e.ID,
		Name:              e.Name,
		EncryptedVaultKey: e.EncryptedVaultKey,
		VaultKeyIV:        e.VaultKeyIv,
		VaultKeyAuthTag:   e.VaultKeyAuthTag,
	}
}

func (r *VaultRepository) FindAll(ctx context.Context, userID string) ([]*domain.Vault, error) {
	vaults, err := r.client.Vault.Query().
		Where(vault.HasUserWith(user.IDEQ(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Vault, len(vaults))
	for i, v := range vaults {
		result[i] = toDomainVault(v)
	}
	return result, nil
}

func (r *VaultRepository) FindByID(ctx context.Context, id string) (*domain.Vault, error) {
	v, err := r.client.Vault.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainVault(v), nil
}

func (r *VaultRepository) Create(ctx context.Context, v *domain.Vault) (*domain.Vault, error) {
	created, err := r.client.Vault.Create().
		SetName(v.Name).
		SetEncryptedVaultKey(v.EncryptedVaultKey).
		SetVaultKeyIv(v.VaultKeyIV).
		SetVaultKeyAuthTag(v.VaultKeyAuthTag).
		SetUserID(v.UserID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainVault(created), nil
}

func (r *VaultRepository) Update(ctx context.Context, v *domain.Vault) (*domain.Vault, error) {
	updated, err := r.client.Vault.UpdateOneID(v.ID).
		SetName(v.Name).
		SetEncryptedVaultKey(v.EncryptedVaultKey).
		SetVaultKeyIv(v.VaultKeyIV).
		SetVaultKeyAuthTag(v.VaultKeyAuthTag).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainVault(updated), nil
}

func (r *VaultRepository) Delete(ctx context.Context, id string) error {
	return r.client.Vault.DeleteOneID(id).Exec(ctx)
}
