package persistence

import (
	"context"

	domain "github.com/dannegm/anubix-server/cmd/internal/domain/device"
	"github.com/dannegm/anubix-server/ent"
	"github.com/dannegm/anubix-server/ent/device"
	"github.com/dannegm/anubix-server/ent/user"
)

type DeviceRepository struct {
	client *ent.Client
}

func NewDeviceRepository(client *ent.Client) *DeviceRepository {
	return &DeviceRepository{client: client}
}

func toDomainDevice(e *ent.Device) *domain.Device {
	return &domain.Device{
		ID:          e.ID,
		Name:        e.Name,
		Fingerprint: e.Fingerprint,
		DeviceType:  e.DeviceType.String(),
	}
}

func (r *DeviceRepository) FindAll(ctx context.Context, userID string) ([]*domain.Device, error) {
	devices, err := r.client.Device.Query().
		Where(device.HasUserWith(user.IDEQ(userID))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Device, len(devices))
	for i, d := range devices {
		result[i] = toDomainDevice(d)
	}
	return result, nil
}

func (r *DeviceRepository) FindByID(ctx context.Context, id string) (*domain.Device, error) {
	d, err := r.client.Device.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainDevice(d), nil
}

func (r *DeviceRepository) Create(ctx context.Context, d *domain.Device) (*domain.Device, error) {
	created, err := r.client.Device.Create().
		SetName(d.Name).
		SetFingerprint(d.Fingerprint).
		SetDeviceType(device.DeviceType(d.DeviceType)).
		SetUserID(d.UserID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainDevice(created), nil
}

func (r *DeviceRepository) Update(ctx context.Context, d *domain.Device) (*domain.Device, error) {
	updated, err := r.client.Device.UpdateOneID(d.ID).
		SetName(d.Name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainDevice(updated), nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	return r.client.Device.DeleteOneID(id).Exec(ctx)
}
