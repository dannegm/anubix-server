package http

import (
	internaljwt "github.com/dannegm/anubix-server/cmd/internal/infrastructure/jwt"
	"github.com/dannegm/anubix-server/cmd/internal/interfaces/http/handlers"
	"github.com/dannegm/anubix-server/cmd/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(
	jwtManager *internaljwt.Manager,
	authHandler *handlers.AuthHandler,
	meHandler *handlers.MeHandler,
	vaultHandler *handlers.VaultHandler,
	deviceHandler *handlers.DeviceHandler,
) *Router {
	r := gin.Default()

	auth := middleware.Auth(jwtManager)
	requireDevice := middleware.RequireDevice(jwtManager)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	authHandler.RegisterRoutes(r, auth)
	meHandler.RegisterRoutes(r, auth)
	vaultHandler.RegisterRoutes(r, requireDevice)
	deviceHandler.RegisterRoutes(r, auth)
	return &Router{engine: r}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
