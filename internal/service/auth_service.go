package service

import (
	"context"
	"database/sql"
	"errors"
	db "jahitin_be/database/repository"
	"jahitin_be/token"
)

type AuthService interface {
	CreateLocalSession(ctx context.Context, deviceID string) (string, *token.Payload, error)
	GetSession(ctx context.Context, user_id int64) bool
	CreateAccountSession(user_id int64)
}

type authService struct {
	store      db.Store
	tokenMaker token.Maker
}

func NewAuthService(store db.Store, tokenMaker token.Maker) AuthService {
	return &authService{store: store, tokenMaker: tokenMaker}
}

func (s *authService) GetSession(ctx context.Context, user_id int64) bool {

	_, err := s.store.GetUserByID(ctx, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		return false
	}

	return true
}

func (s *authService) CreateLocalSession(ctx context.Context, deviceID string) (string, *token.Payload, error) {
	if deviceID == "" {
		return "", nil, errors.New("device id is required")
	}

	// get user by device id
	user, err := s.store.GetUserByDeviceID(ctx, sql.NullString{String: deviceID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}

	tokenValue, tokenPayload, err := s.tokenMaker.CreateToken(
		user.ID,
		user.Name,
		token.AuthTypeLocal,
		deviceID,
		user.Username.String,
		user.Email.String,
	)
	if err != nil {
		return "", nil, err
	}

	return tokenValue, tokenPayload, nil

}

func (s *authService) CreateAccountSession(user_id int64) {

}
