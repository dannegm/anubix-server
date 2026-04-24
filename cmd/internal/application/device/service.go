package device

import (
	"context"

	domain "github.com/dannegm/anubix-server/cmd/internal/domain/device"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(ctx context.Context, userID string) ([]*domain.Device, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, d *domain.Device) (*domain.Device, error) {
	return s.repo.Create(ctx, d)
}

func (s *Service) Update(ctx context.Context, d *domain.Device) (*domain.Device, error) {
	return s.repo.Update(ctx, d)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
