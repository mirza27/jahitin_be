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
	authRoutes.GET("/category/list", h.ClothesCategory.ListAllClothesCategoriesHandler)

	// services
	authRoutes.GET("/service/list", h.ServiceType.ListAllServiceTypesHandler)

	// order
	authRoutes.POST("/order/create", h.Order.CreateOrderHandler)
	authRoutes.GET("/order/list", h.Order.ListOrdersHandler)
	authRoutes.GET("/order/detail/:order_id", h.Order.GetDetailOrderHandler)

	authRoutes.PUT("/order/update/:order_id", h.Order.UpdateOrderHandler)
	authRoutes.PUT("/order/update/:order_id/items", h.Order.UpdateOrderItemsHandler)
	authRoutes.PUT("/order/status/update", h.Order.UpdateOrderStatusHandler)
	authRoutes.DELETE("/order/delete", h.Order.DeleteOrderHandler)

	// customers
	authRoutes.GET("/customer/list", h.Customer.ListCustomersHandler)
	authRoutes.GET("/customer/detail", h.Customer.GetCustomerDetailHandler)

	s.router = router
}
