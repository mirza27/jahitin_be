package handler

import (
	"database/sql"
	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"
)

type Handler struct {
	User            *UserHandler
	Auth            *AuthHandler
	Order           *OrderHandler
	Customer        *CustomerHandler
	ClothesCategory *ClothesCategoryHandler
	ServiceType     *ServiceTypeHandler
}

func NewHandler(store db.Store, token token.Maker) *Handler {
	userService := service.NewUserService(store)
	authService := service.NewAuthService(store, token)
	orderService := service.NewOrderService(store)
	customerService := service.NewCustomerService(store)
	clothesCategoryService := service.NewClothesCategoryService(store)
	serviceTypeService := service.NewServiceTypeService(store)

	return &Handler{
		User:            NewUserHandler(store, userService),
		Auth:            NewAuthHandler(store, authService, token),
		Order:           NewOrderHandler(store, orderService),
		Customer:        NewCustomerHandler(store, customerService),
		ClothesCategory: NewClothesCategoryHandler(store, clothesCategoryService),
		ServiceType:     NewServiceTypeHandler(store, serviceTypeService),
	}
}

// helper function
func NullStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}
