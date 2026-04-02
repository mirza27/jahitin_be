package api

import (
	"jahitin_be/internal/handler"

	"github.com/gin-gonic/gin"
)

func (s *Server) SetupRoutes() {

	router := gin.Default()
	h := handler.NewHandler(s.store)

	// auth
	router.GET("/auth/session", h.Auth.Session)
	router.POST("/auth/session/create", h.Auth.CreateSession)
	router.POST("/auth/login", h.Auth.Login)

	// categories
	router.GET("/category/list", h.Order.ListCategories)

	// services
	router.GET("/service/list", h.Order.ListServices)

	// middleware for session
	authRoutes := router.Group("/").Use()

	// user
	authRoutes.POST("/user/register/account", h.User.CreateUserAccount)
	authRoutes.POST("/user/register/local", h.User.CreateUserLocal)

	// order
	authRoutes.POST("/order/create", h.Order.CreateOrder)
	authRoutes.GET("/order/list", h.Order.ListOrders)
	authRoutes.GET("/order/detail", h.Order.DetailOrder)
	authRoutes.PUT("/order/update", h.Order.UpdateOrder)
	authRoutes.DELETE("/order/delete", h.Order.DeleteOrder)

	// customers
	authRoutes.GET("/customer/list", h.User.ListCustomers)

	s.router = router
}
