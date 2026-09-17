package handler

import (
	"net/http"

	db "jahitin_be/database/repository"
	"jahitin_be/internal/service"

	"github.com/gin-gonic/gin"
)

type ClothesCategoryHandler struct {
	store   db.Store
	service service.ClothesCategoryService
}

func NewClothesCategoryHandler(store db.Store, categoryService service.ClothesCategoryService) *ClothesCategoryHandler {
	return &ClothesCategoryHandler{store: store, service: categoryService}
}

func (h *ClothesCategoryHandler) ListAllClothesCategoriesHandler(c *gin.Context) {

	categories, err := h.service.ListClothesCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to list clothes categories",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success get clothes categories",
		"data":    categories,
	})

}
