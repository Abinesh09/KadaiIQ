package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

type JWTClaims struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	TenantID     uuid.UUID `json:"tenant_id,omitempty"`
	RoleID       uuid.UUID `json:"role_id,omitempty"`
	RoleName     string    `json:"role_name,omitempty"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	Permissions  []string  `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(
	userID uuid.UUID,
	email string,
	tenantID uuid.UUID,
	roleID uuid.UUID,
	roleName string,
	isSuperAdmin bool,
	permissions []string,
	secret string,
	duration time.Duration,
) (string, error) {
	claims := JWTClaims{
		UserID:       userID,
		Email:        email,
		TenantID:     tenantID,
		RoleID:       roleID,
		RoleName:     roleName,
		IsSuperAdmin: isSuperAdmin,
		Permissions:  permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "KadaiIQ",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenStr string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
