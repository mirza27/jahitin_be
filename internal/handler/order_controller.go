package handler

import (
	"net/http"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	store   db.Store
	service service.OrderService
}

func NewOrderHandler(store db.Store, orderService service.OrderService) *OrderHandler {
	return &OrderHandler{store: store, service: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) DetailOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) ListCategories(c *gin.Context) {
	categories, err := h.store.ListClothesCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *OrderHandler) ListServices(c *gin.Context) {
	serviceTypes, err := h.store.ListServiceTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch service types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceTypes})
}
