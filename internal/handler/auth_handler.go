package handler

import (
	"net/http"
	"strconv"
	"time"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	store      db.Store
	service    service.AuthService
	tokenMaker token.Maker
}

func NewAuthHandler(store db.Store, authService service.AuthService, tokenMaker token.Maker) *AuthHandler {
	return &AuthHandler{store: store, service: authService, tokenMaker: tokenMaker}
}

type SessionResponse struct {
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	UserType  string    `json:"user_type"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (h *AuthHandler) Session(c *gin.Context) {

	authPayload := c.MustGet("authorization_payload").(*token.Payload)

	isValid := h.service.GetSession(c.Request.Context(), authPayload.UserID)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	response := SessionResponse{
		UserId:    strconv.FormatInt(authPayload.UserID, 10),
		Name:      authPayload.Name,
		Username:  authPayload.Username,
		Email:     authPayload.Email,
		ExpiredAt: time.Unix(authPayload.ExpiredAt, 0),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

type LoginLocalRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}
type LoginResponse struct {
	Token string `json:"token"`
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

	c.JSON(http.StatusOK, gin.H{"data": LoginResponse{Token: token}})

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
