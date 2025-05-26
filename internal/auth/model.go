package auth

import (
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/google/uuid"
)

type AuthenticationResult struct {
	ID           uuid.UUID
	AccessToken  *jwt.AccessToken
	RefreshToken *jwt.RefreshToken
}
