package auth

import (
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/pkg/crypto"
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r *RegisterRequest) ToRegisterParams() (*repository.RegisterParams, error) {
	passHash, err := crypto.HashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	return &repository.RegisterParams{
		ID:           uuid.New(),
		Username:     r.Username,
		PasswordHash: passHash,
	}, nil
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthenticationResponse struct {
	ID          uuid.UUID        `json:"id,omitempty"`
	AccessToken *jwt.AccessToken `json:"accessToken"`
}
