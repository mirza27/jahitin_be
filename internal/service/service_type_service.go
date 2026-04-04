package service

import (
	"context"
	db "jahitin_be/database/repository"
)

type ServiceTypeService interface {
	ListServiceTypes(ctx context.Context) ([]*db.ServiceType, error)
}

type serviceTypeService struct {
	store db.Store
}

func NewServiceTypeService(store db.Store) ServiceTypeService {
	return &serviceTypeService{store: store}
}

func (s *serviceTypeService) ListServiceTypes(ctx context.Context) ([]*db.ServiceType, error) {

	services, err := s.store.ListServiceTypes(ctx)
	if err != nil {
		return nil, err
	}

	servicePtrs := make([]*db.ServiceType, len(services))
	for i, service := range services {
		servicePtrs[i] = &service
	}

	return servicePtrs, nil

}
