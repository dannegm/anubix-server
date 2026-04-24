package vault

import (
	"context"

	domain "github.com/dannegm/anubix-server/cmd/internal/domain/vault"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, userID string) ([]*domain.Vault, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Vault, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, v *domain.Vault) (*domain.Vault, error) {
	return s.repo.Create(ctx, v)
}

func (s *Service) Update(ctx context.Context, v *domain.Vault) (*domain.Vault, error) {
	return s.repo.Update(ctx, v)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
