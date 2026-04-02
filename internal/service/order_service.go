package service

import db "jahitin_be/database/repository"

type OrderService interface{}

type orderService struct {
	store db.Store
}

func NewOrderService(store db.Store) OrderService {
	return &orderService{store: store}
}
