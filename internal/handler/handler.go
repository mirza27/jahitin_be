package handler

import (
	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
)

type Handler struct {
	User  *UserHandler
	Auth  *AuthHandler
	Order *OrderHandler
}

func NewHandler(store db.Store) *Handler {
	userService := service.NewUserService(store)
	authService := service.NewAuthService(store)
	orderService := service.NewOrderService(store)

	return &Handler{
		User:  NewUserHandler(store, userService),
		Auth:  NewAuthHandler(store, authService),
		Order: NewOrderHandler(store, orderService),
	}
}
