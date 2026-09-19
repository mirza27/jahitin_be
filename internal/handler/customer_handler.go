package handler

import (
	"net/http"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"
	"jahitin_be/token"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	store   db.Store
	service service.CustomerService
}

func NewCustomerHandler(store db.Store, customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{store: store, service: customerService}
}

type ListCustomerListRequest struct {
	Search string `form:"search"`
	Page   int32  `form:"page" binding:"required,min=1"`
	Limit  int32  `form:"limit" binding:"required,min=1,max=100"`
}

type ListCustomerResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (h *CustomerHandler) ListCustomersHandler(c *gin.Context) {
	var req ListCustomerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := c.MustGet("authorization_payload").(*token.Payload)
	userID := authPayload.UserID

	customerFilter := service.ListCustomersFilter{
		UserID: userID,
		Search: req.Search,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	customers, err := h.service.ListUserCustomers(c.Request.Context(), customerFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "failed to list customers",
				"details": err.Error()})
		return
	}

	customerListResponse := make([]ListCustomerResponse, len(customers))
	for i, customer := range customers {
		customerListResponse[i] = ListCustomerResponse{
			ID:        customer.ID,
			Name:      customer.Name,
			Phone:     NullStringValue(customer.Phone),
			Notes:     NullStringValue(customer.Notes),
			CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": customerListResponse})
}

type GetCustomerDetailRequest struct {
	ID int64 `form:"id" binding:"required"`
}

func (h *CustomerHandler) GetCustomerDetailHandler(c *gin.Context) {
	var req GetCustomerDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := c.MustGet("authorization_payload").(*token.Payload)
	userID := authPayload.UserID

	customer, err := h.service.GetUserCustomerDetail(c.Request.Context(), userID, req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get customer detail", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerHandler) GetDetailCustomer(c *gin.Context) {

}
