package handler

import (
	"database/sql"
	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"
)

type Handler struct {
	User     *UserHandler
	Auth     *AuthHandler
	Order    *OrderHandler
	Customer *CustomerHandler
}

func NewHandler(store db.Store, token token.Maker) *Handler {
	userService := service.NewUserService(store)
	authService := service.NewAuthService(store, token)
	orderService := service.NewOrderService(store)
	customerService := service.NewCustomerService(store)

	return &Handler{
		User:     NewUserHandler(store, userService),
		Auth:     NewAuthHandler(store, authService, token),
		Order:    NewOrderHandler(store, orderService),
		Customer: NewCustomerHandler(store, customerService),
	}
}

// helper function
func NullStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
