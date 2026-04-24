package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	userapp "github.com/dannegm/anubix-server/cmd/internal/application/user"
	persistence "github.com/dannegm/anubix-server/cmd/internal/infrastructure/persistence/ent"
	userhttp "github.com/dannegm/anubix-server/cmd/internal/interfaces/http"
	"github.com/dannegm/anubix-server/ent"
)

func main() {
	_ = godotenv.Load()

	client, err := ent.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Wiring
	userRepo := persistence.NewUserRepository(client)
	userService := userapp.NewService(userRepo)
	userHandler := userhttp.NewUserHandler(userService)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	userHandler.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	r.Run(":" + port)
}
