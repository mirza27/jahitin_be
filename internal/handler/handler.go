package handler

import (
	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"
)

type Handler struct {
	User  *UserHandler
	Auth  *AuthHandler
	Order *OrderHandler
}

func NewHandler(store db.Store, token token.Maker) *Handler {
	userService := service.NewUserService(store)
	authService := service.NewAuthService(store, token)
	orderService := service.NewOrderService(store)

	return &Handler{
		User:  NewUserHandler(store, userService),
		Auth:  NewAuthHandler(store, authService, token),
		Order: NewOrderHandler(store, orderService),
	}
}
