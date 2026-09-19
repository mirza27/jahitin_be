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
	DeviceID string `json:"device_id" binding:"required"`
}

type CreateUserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	DeviceID string `json:"device_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	UserType string `json:"user_type"`
}

func (h *UserHandler) CreateUserLocal(c *gin.Context) {
	var req CreateUserLocalRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// create base user
	user, err := h.service.CreateNewLocal(c.Request.Context(), service.CreateLocalUserInput{
		Name:     req.Name,
		DeviceID: req.DeviceID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create user",
			"error":   err.Error(),
		})
		return
	}

	userResponse := CreateUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		DeviceID: user.DeviceID.String,
		UserType: user.UserType,
		Username: NullStringValue(user.Username),
		Email:    NullStringValue(user.Email),
		Phone:    NullStringValue(user.Phone),
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "user created successfully",
		"data":    userResponse,
	})
}

func (h *UserHandler) CreateUserAccount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *UserHandler) ListCustomers(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})

}
