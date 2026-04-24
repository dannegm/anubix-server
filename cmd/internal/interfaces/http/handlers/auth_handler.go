package handlers

import (
	"net/http"

	appauth "github.com/dannegm/anubix-server/cmd/internal/application/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *appauth.Service
}

func NewAuthHandler(service *appauth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	g := r.Group("/auth")
	g.POST("/register", h.Register)
	g.GET("/salt", h.GetSalt)
	g.POST("/login", h.Login)
	g.POST("/token", auth, h.Token)
}

type registerRequest struct {
	Email             string `json:"email" binding:"required,email"`
	AuthHash          string `json:"auth_hash" binding:"required"`
	Salt              string `json:"salt" binding:"required"`
	VaultName         string `json:"vault_name" binding:"required"`
	EncryptedVaultKey string `json:"encrypted_vault_key" binding:"required"`
	VaultKeyIV        string `json:"vault_key_iv" binding:"required"`
	VaultKeyAuthTag   string `json:"vault_key_auth_tag" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	AuthHash string `json:"auth_hash" binding:"required"`
}

type tokenRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Register(c.Request.Context(), appauth.RegisterInput{
		Email:             req.Email,
		AuthHash:          req.AuthHash,
		Salt:              req.Salt,
		VaultName:         req.VaultName,
		EncryptedVaultKey: req.EncryptedVaultKey,
		VaultKeyIV:        req.VaultKeyIV,
		VaultKeyAuthTag:   req.VaultKeyAuthTag,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) GetSalt(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	salt, err := h.service.GetSalt(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"salt": salt})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), appauth.LoginInput{
		Email:    req.Email,
		AuthHash: req.AuthHash,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Token(c *gin.Context) {
	userID := c.GetString("user_id")

	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Token(c.Request.Context(), userID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
