package handlers

import (
	"net/http"

	appdevice "github.com/dannegm/anubix-server/cmd/internal/application/device"
	domain "github.com/dannegm/anubix-server/cmd/internal/domain/device"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	service *appdevice.Service
}

func NewDeviceHandler(service *appdevice.Service) *DeviceHandler {
	return &DeviceHandler{service: service}
}

type deviceRequest struct {
	Name        string `json:"name" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
	DeviceType  string `json:"device_type" binding:"required,oneof=web ios android desktop"`
}

func (h *DeviceHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	g := r.Group("/devices", auth)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *DeviceHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")
	devices, err := h.service.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) GetByID(c *gin.Context) {
	d, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DeviceHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req deviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.service.Create(c.Request.Context(), &domain.Device{
		UserID:      userID,
		Name:        req.Name,
		Fingerprint: req.Fingerprint,
		DeviceType:  req.DeviceType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DeviceHandler) Update(c *gin.Context) {
	var req deviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d, err := h.service.Update(c.Request.Context(), &domain.Device{
		ID:   c.Param("id"),
		Name: req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
