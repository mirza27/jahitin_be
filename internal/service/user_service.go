package service

import (
	"context"
	"database/sql"

	db "jahitin_be/database/repository"
)

type UserService interface {
	CreateNewLocal(ctx context.Context, input CreateLocalUserInput) (db.User, error)
}

type CreateLocalUserInput struct {
	Name     string
	DeviceID string
}

type userService struct {
	store db.Store
}

func NewUserService(store db.Store) UserService {
	return &userService{store: store}
}

func (s *userService) CreateNewLocal(ctx context.Context, input CreateLocalUserInput) (db.User, error) {

	// create base user
	params := db.CreateUserParams{
		Name:         input.Name,
		Username:     sql.NullString{Valid: false},
		Email:        sql.NullString{Valid: false},
		PasswordHash: sql.NullString{Valid: false},
		DeviceID:     sql.NullString{String: input.DeviceID, Valid: true},
		Phone:        sql.NullString{Valid: false},
		UserType:     "tailor",
	}

	return s.store.CreateUser(ctx, params)
}
