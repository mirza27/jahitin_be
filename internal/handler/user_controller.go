package handler

import (
	"net/http"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store   db.Store
	service service.UserService
}

func NewUserHandler(store db.Store, userService service.UserService) *UserHandler {
	return &UserHandler{store: store, service: userService}
}

type CreateUserLocalRequest struct {
	Name     string `json:"name" binding:"required"`
	DeviceID string `json:"device_id" binding:"required,string"`
}

func (h *UserHandler) CreateUserLocal(c *gin.Context) {
	var req CreateUserLocalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create base user
	user, err := h.service.CreateNewLocal(c.Request.Context(), service.CreateLocalUserInput{
		Name:     req.Name,
		DeviceID: req.DeviceID,
	})

	// generate token
	// token, err

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) CreateUserAccount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *UserHandler) ListCustomers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})

}
