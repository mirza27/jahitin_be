package handler

import (
	"net/http"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	store   db.Store
	service service.AuthService
}

func NewAuthHandler(store db.Store, authService service.AuthService) *AuthHandler {
	return &AuthHandler{store: store, service: authService}
}

func (h *AuthHandler) Session(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

type LoginLocalRequest struct {
	DeviceID string `json:"device_id" binding:"required,string"`
}

func (h *AuthHandler) LoginLocal(c *gin.Context) {
	var req LoginLocalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, _, err := h.service.CreateLocalSession(c.Request.Context(), req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

type LoginAccountRequest struct {
	Phone    string `json:"phone" binding:"optional"`
	Username string `json:"username" binding:"optional"`
	Email    string `json:"email" binding:"optional"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) LoginAccount(c *gin.Context) {
	var req LoginAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// atleast one of phone, username, email must be provided
	if req.Phone == "" && req.Username == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one of phone, username, email must be provided"})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
