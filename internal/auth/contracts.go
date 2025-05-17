package auth

import (
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	ID          uuid.UUID        `json:"id"`
	AccessToken *jwt.AccessToken `json:"accessToken"`
}
