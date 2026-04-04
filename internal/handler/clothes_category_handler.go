package handler

import (
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
		c.JSON(500, gin.H{"error": "failed to list clothes categories"})
		return
	}

	c.JSON(200, gin.H{"data": categories})

}
