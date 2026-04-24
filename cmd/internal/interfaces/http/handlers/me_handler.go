package handlers

import (
	"log"
	"net/http"

	domainuser "github.com/dannegm/anubix-server/cmd/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	userRepo domainuser.Repository
}

func NewMeHandler(userRepo domainuser.Repository) *MeHandler {
	return &MeHandler{userRepo: userRepo}
}

func (h *MeHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	g := r.Group("/me", auth)
	g.GET("", h.Me)
}

func (h *MeHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	log.Printf("user_id: %s", userID)

	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
