package handler

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to list service types",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success get service types",
		"data":    services,
	})

}
