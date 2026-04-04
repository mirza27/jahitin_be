package handler

import (
	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type ServiceTypeHandler struct {
	store   db.Store
	service service.ServiceTypeService
}

func NewServiceTypeHandler(store db.Store, service service.ServiceTypeService) *ServiceTypeHandler {
	return &ServiceTypeHandler{store: store, service: service}
}

func (h *ServiceTypeHandler) ListAllServiceTypesHandler(c *gin.Context) {

	services, err := h.service.ListServiceTypes(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to list service types"})
		return
	}

	c.JSON(200, gin.H{"data": services})

}
