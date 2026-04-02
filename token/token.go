package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	AuthTypeLocal        = "local"
	AuthTypeAccount      = "account"
	LocalTokenDuration   = 24 * time.Hour
	AccountTokenDuration = 7 * 24 * time.Hour
)

var (
	ErrInvalidAuthType = errors.New("invalid auth_type")
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token is expired")
)

type Payload struct {
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	AuthType  string `json:"auth_type"`
	DeviceID  string `json:"device_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	ExpiredAt int64  `json:"expired_at"`
}

type Maker struct {
	secretKey []byte
}

func NewMaker(secret string) (*Maker, error) {
	if len(secret) < 32 {
		return nil, errors.New("secret key must be at least 32 characters")
	}

	return &Maker{secretKey: []byte(secret)}, nil
}

func NewPayload(userID int64, name string, authType string, deviceID string, username string, email string, expiredAt int64) (*Payload, error) {
	if !isValidAuthType(authType) {
		return nil, ErrInvalidAuthType
	}

	return &Payload{
		UserID:    userID,
		Name:      name,
		AuthType:  authType,
		DeviceID:  deviceID,
		Username:  username,
		Email:     email,
		ExpiredAt: expiredAt,
	}, nil
}

func (m *Maker) CreateToken(userID int64, name string, authType string, deviceID string, username string, email string) (string, *Payload, error) {
	var duration time.Duration
	switch authType {
	case AuthTypeLocal:
		duration = LocalTokenDuration
	case AuthTypeAccount:
		duration = AccountTokenDuration
	default:
		return "", nil, ErrInvalidAuthType
	}

	expiredAt := time.Now().Add(duration).Unix()

	payload, err := NewPayload(userID, name, authType, deviceID, username, email, expiredAt)
	if err != nil {
		return "", nil, err
	}

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	payloadPart := base64.RawURLEncoding.EncodeToString(payloadBytes)
	unsigned := header + "." + payloadPart
	signature := m.sign(unsigned)

	token := unsigned + "." + signature
	return token, payload, nil
}

func (m *Maker) VerifyToken(token string) (*Payload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expected := m.sign(unsigned)
	if subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expected)) != 1 {
		return nil, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var payload Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidToken
	}

	if !isValidAuthType(payload.AuthType) {
		return nil, ErrInvalidAuthType
	}

	if payload.ExpiredAt <= time.Now().Unix() {
		return nil, ErrExpiredToken
	}

	return &payload, nil
}

func (m *Maker) sign(data string) string {
	h := hmac.New(sha256.New, m.secretKey)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func isValidAuthType(authType string) bool {
	switch authType {
	case AuthTypeLocal, AuthTypeAccount:
		return true
	default:
		return false
	}
}

func AuthTypeFromString(val string) (string, error) {
	if !isValidAuthType(val) {
		return "", fmt.Errorf("%w: %s", ErrInvalidAuthType, val)
	}

	return val, nil
}
