package device

import "context"

type Repository interface {
	FindAll(ctx context.Context, userID string) ([]*Device, error)
	FindByID(ctx context.Context, id string) (*Device, error)
	Create(ctx context.Context, d *Device) (*Device, error)
	Update(ctx context.Context, d *Device) (*Device, error)
	Delete(ctx context.Context, id string) error
}
