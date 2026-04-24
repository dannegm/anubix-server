package auth

import (
	"context"
	"errors"

	domainuser "github.com/dannegm/anubix-server/cmd/internal/domain/user"
	domainvault "github.com/dannegm/anubix-server/cmd/internal/domain/vault"
	internaljwt "github.com/dannegm/anubix-server/cmd/internal/infrastructure/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   domainuser.Repository
	vaultRepo  domainvault.Repository
	jwtManager *internaljwt.Manager
}

func NewService(userRepo domainuser.Repository, vaultRepo domainvault.Repository, jwtManager *internaljwt.Manager) *Service {
	return &Service{
		userRepo:   userRepo,
		vaultRepo:  vaultRepo,
		jwtManager: jwtManager,
	}
}

type RegisterInput struct {
	Email             string
	AuthHash          string
	Salt              string
	VaultName         string
	EncryptedVaultKey string
	VaultKeyIV        string
	VaultKeyAuthTag   string
}

type LoginInput struct {
	Email    string
	AuthHash string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (string, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return "", errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.AuthHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.Create(ctx, &domainuser.User{
		Email:    input.Email,
		AuthHash: string(hashed),
		Salt:     input.Salt,
	})
	if err != nil {
		return "", err
	}

	_, err = s.vaultRepo.Create(ctx, &domainvault.Vault{
		UserID:            user.ID,
		Name:              input.VaultName,
		EncryptedVaultKey: input.EncryptedVaultKey,
		VaultKeyIV:        input.VaultKeyIV,
		VaultKeyAuthTag:   input.VaultKeyAuthTag,
	})
	if err != nil {
		return "", err
	}

	return s.jwtManager.Generate(user.ID, nil)
}

func (s *Service) GetSalt(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}
	return user.Salt, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.AuthHash), []byte(input.AuthHash)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwtManager.Generate(user.ID, nil)
}

func (s *Service) Token(ctx context.Context, userID string, deviceID string) (string, error) {
	return s.jwtManager.Generate(userID, &deviceID)
}
