package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/dannegm/anubix-server/cmd/internal/infrastructure/config"
	internaljwt "github.com/dannegm/anubix-server/cmd/internal/infrastructure/jwt"
	persistence "github.com/dannegm/anubix-server/cmd/internal/infrastructure/persistence/ent"
	"github.com/dannegm/anubix-server/ent"

	appauth "github.com/dannegm/anubix-server/cmd/internal/application/auth"
	appdevice "github.com/dannegm/anubix-server/cmd/internal/application/device"
	appvault "github.com/dannegm/anubix-server/cmd/internal/application/vault"
	httprouter "github.com/dannegm/anubix-server/cmd/internal/interfaces/http"
	"github.com/dannegm/anubix-server/cmd/internal/interfaces/http/handlers"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	client, err := ent.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Infrastructure
	jwtManager := internaljwt.NewManager(cfg.JWTSecret)
	userRepo := persistence.NewUserRepository(client)
	vaultRepo := persistence.NewVaultRepository(client)
	deviceRepo := persistence.NewDeviceRepository(client)

	// Services
	authService := appauth.NewService(userRepo, vaultRepo, jwtManager)
	vaultService := appvault.NewService(vaultRepo)
	deviceService := appdevice.NewService(deviceRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	meHandler := handlers.NewMeHandler(userRepo)
	vaultHandler := handlers.NewVaultHandler(vaultService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// Router
	router := httprouter.NewRouter(jwtManager, authHandler, meHandler, vaultHandler, deviceHandler)

	log.Printf("Server running on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
