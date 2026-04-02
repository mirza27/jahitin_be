package service

import db "jahitin_be/database/repository"

type AuthService interface{}

type authService struct {
	store db.Store
}

func NewAuthService(store db.Store) AuthService {
	return &authService{store: store}
}
