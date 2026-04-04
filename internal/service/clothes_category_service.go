package service

import (
	"context"
	db "jahitin_be/database/repository"
)

type ClothesCategoryService interface {
	ListClothesCategories(ctx context.Context) ([]*db.ClothesCategory, error)
}

type clothesCategoryService struct {
	store db.Store
}

func NewClothesCategoryService(store db.Store) ClothesCategoryService {
	return &clothesCategoryService{store: store}
}

func (s *clothesCategoryService) ListClothesCategories(ctx context.Context) ([]*db.ClothesCategory, error) {
	categories, err := s.store.ListClothesCategories(ctx)
	if err != nil {
		return nil, err
	}

	categoryPtrs := make([]*db.ClothesCategory, len(categories))
	for i := range categories {
		categoryPtrs[i] = &categories[i]
	}

	return categoryPtrs, nil
}
