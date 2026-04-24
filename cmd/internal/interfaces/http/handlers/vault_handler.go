package handlers

import (
	"log"
	"net/http"

	appvault "github.com/dannegm/anubix-server/cmd/internal/application/vault"
	domainvault "github.com/dannegm/anubix-server/cmd/internal/domain/vault"
	"github.com/gin-gonic/gin"
)

type VaultHandler struct {
	service *appvault.Service
}

func NewVaultHandler(service *appvault.Service) *VaultHandler {
	return &VaultHandler{service: service}
}

func (h *VaultHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	g := r.Group("/vaults", auth)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

type vaultRequest struct {
	Name              string `json:"name" binding:"required"`
	EncryptedVaultKey string `json:"encrypted_vault_key" binding:"required"`
	VaultKeyIV        string `json:"vault_key_iv" binding:"required"`
	VaultKeyAuthTag   string `json:"vault_key_auth_tag" binding:"required"`
}

func (h *VaultHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")
	vaults, err := h.service.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vaults)
}

func (h *VaultHandler) GetByID(c *gin.Context) {
	vault, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "vault not found"})
		return
	}
	c.JSON(http.StatusOK, vault)
}

func (h *VaultHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	log.Printf("user_id: %s", userID)

	var req vaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vault, err := h.service.Create(c.Request.Context(), &domainvault.Vault{
		UserID:            userID,
		Name:              req.Name,
		EncryptedVaultKey: req.EncryptedVaultKey,
		VaultKeyIV:        req.VaultKeyIV,
		VaultKeyAuthTag:   req.VaultKeyAuthTag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vault)
}

func (h *VaultHandler) Update(c *gin.Context) {
	var req vaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vault, err := h.service.Update(c.Request.Context(), &domainvault.Vault{
		ID:                c.Param("id"),
		Name:              req.Name,
		EncryptedVaultKey: req.EncryptedVaultKey,
		VaultKeyIV:        req.VaultKeyIV,
		VaultKeyAuthTag:   req.VaultKeyAuthTag,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vault)
}

func (h *VaultHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
