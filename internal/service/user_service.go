package service

import (
	"context"
	"database/sql"

	db "jahitin_be/database/repository"
)

type UserService interface {
	CreateLocal(ctx context.Context, input CreateLocalUserInput) (db.User, error)
}

type CreateLocalUserInput struct {
	Name     string
	Username string
	Email    string
	Password string
	DeviceID string
	Phone    string
	UserType string
}

type userService struct {
	store db.Store
}

func NewUserService(store db.Store) UserService {
	return &userService{store: store}
}

func (s *userService) CreateLocal(ctx context.Context, input CreateLocalUserInput) (db.User, error) {
	if input.UserType == "" {
		input.UserType = "tailor"
	}

	params := db.CreateUserParams{
		Name:         input.Name,
		Username:     toNullString(input.Username),
		Email:        toNullString(input.Email),
		PasswordHash: toNullString(input.Password),
		DeviceID:     toNullString(input.DeviceID),
		Phone:        toNullString(input.Phone),
		UserType:     input.UserType,
	}

	return s.store.CreateUser(ctx, params)
}

func toNullString(val string) sql.NullString {
	if val == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: val, Valid: true}
}
