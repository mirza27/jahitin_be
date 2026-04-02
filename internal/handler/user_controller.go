package handler

import (
	"net/http"
	"strconv"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store   db.Store
	service service.UserService
}

type createUserLocalRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	UserType string `json:"user_type"`
}

func NewUserHandler(store db.Store, userService service.UserService) *UserHandler {
	return &UserHandler{store: store, service: userService}
}

func (h *UserHandler) CreateUserLocal(c *gin.Context) {
	var req createUserLocalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateLocal(c.Request.Context(), service.CreateLocalUserInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		DeviceID: req.DeviceID,
		Phone:    req.Phone,
		UserType: req.UserType,
	})
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
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	customers, err := h.store.ListCustomersByUserID(c.Request.Context(), db.ListCustomersByUserIDParams{
		UserID: userID,
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})

}
