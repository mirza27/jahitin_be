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

func (h *AuthHandler) CreateSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}
