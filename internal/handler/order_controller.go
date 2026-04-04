package handler

import (
	"net/http"
	"time"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	store   db.Store
	service service.OrderService
}

func NewOrderHandler(store db.Store, orderService service.OrderService) *OrderHandler {
	return &OrderHandler{store: store, service: orderService}
}

type CustomerOrderDetailRequest struct {
	CustomerID    *int64 `json:"customer_id" binding:"required_if=IsNewCustomer false"`
	CustomerName  string `json:"customer_name" binding:"required_if=IsNewCustomer true"`
	CustomerPhone string `json:"customer_phone" binding:"required_if=IsNewCustomer true"`
	IsNewCustomer *bool  `json:"is_new_customer" binding:"required"`
}

type OrderItemsDetailRequest struct {
	ClothesFor          string `json:"clothes_for" binding:"required"`
	Notes               string `json:"notes"`
	ClothesCategoryID   *int64 `json:"clothes_category_id"`
	ServiceTypeID       *int64 `json:"service_type_id" binding:"omitempty,gt=0,excluded_with=CustomServiceName"`
	CustomServiceName   string `json:"custom_service_name" binding:"required_without=ServiceTypeID,excluded_with=ServiceTypeID"`
	Price               *int64 `json:"price" binding:"required,min=0"`
	IsSaveCustomerNotes *bool  `json:"is_save_customer_notes" binding:"required"`
}

type CreateOrderRequest struct {
	Name       string                     `json:"name" binding:"required"`
	Deadline   string                     `json:"deadline" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	Customer   CustomerOrderDetailRequest `json:"customer" binding:"required"`
	OrderItems []OrderItemsDetailRequest  `json:"order_items" binding:"required"`
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := c.MustGet("authorization_payload").(*token.Payload)
	userID := authPayload.UserID

	// check format deadline
	var deadline *time.Time
	if req.Deadline != "" {
		parsedDeadline, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid deadline format, expected RFC3339"})
			return
		}
		deadline = &parsedDeadline
	}

	// build order items
	items := make([]service.OrderItemInput, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		items = append(items, service.OrderItemInput{
			ClothesFor:          item.ClothesFor,
			Notes:               item.Notes,
			ClothesCategoryID:   *item.ClothesCategoryID,
			ServiceTypeID:       *item.ServiceTypeID,
			CustomServiceName:   item.CustomServiceName,
			Price:               *item.Price,
			IsSaveCustomerNotes: *item.IsSaveCustomerNotes, // get bool val
		})
	}

	oInput := service.CreateUserOrderInput{
		Name:     req.Name,
		Deadline: deadline,
		Items:    items,
	}

	cInput := service.CustomerInput{
		IsNewCustomer: *req.Customer.IsNewCustomer, // get bool val
		CustomerID:    *req.Customer.CustomerID,
		Name:          req.Customer.CustomerName,
		Phone:         req.Customer.CustomerPhone,
	}

	valid, err := h.service.CreateUserOrder(c.Request.Context(), userID, oInput, cInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create order", "details": err.Error()})
		return
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order data"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "order created successfully"})
}

type ListOrdersRequest struct {
	Status  string `form:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Search  string `form:"search" binding:"omitempty"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=created_at deadline"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=asc desc"`
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	var req ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{"message": "not implemented"})
}

func (h *OrderHandler) GetRelatedCustomerOrderHandler(c *gin.Context) {

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
