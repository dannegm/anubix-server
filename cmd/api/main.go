package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carga .env si existe; si no, usa las env vars del sistema (Docker/DO)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello Gopher!"})
	})

	log.Printf("Server running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
