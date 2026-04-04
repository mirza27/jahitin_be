package api

import (
	"jahitin_be/internal/handler"

	"github.com/gin-gonic/gin"
)

func (s *ApiServer) SetupRoutes() {

	router := gin.Default()
	h := handler.NewHandler(s.store, s.tokenmaker)
	authRoutes := router.Group("/", s.AuthMiddleware())

	// auth
	router.POST("/auth/login/local", h.Auth.LoginLocal)
	router.POST("/auth/login/account", h.Auth.LoginAccount)

	// user
	router.POST("/user/register/local", h.User.CreateUserLocal)
	router.POST("/user/register/account", h.User.CreateUserAccount)

	// session
	authRoutes.GET("/auth/session", h.Auth.Session)

	// categories
	router.GET("/category/list", h.Order.ListCategories)

	// services
	router.GET("/service/list", h.Order.ListServices)

	// order
	authRoutes.POST("/order/create", h.Order.CreateOrderHandler)
	authRoutes.GET("/order/customer/detail", h.Order.GetRelatedCustomerOrderHandler)
	authRoutes.GET("/order/list", h.Order.ListOrders)
	authRoutes.GET("/order/detail", h.Order.DetailOrder)
	authRoutes.PUT("/order/update", h.Order.UpdateOrder)
	authRoutes.DELETE("/order/delete", h.Order.DeleteOrder)

	// customers
	authRoutes.GET("/customer/list", h.User.ListCustomers)

	s.router = router
}
