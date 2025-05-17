package auth

import (
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/V2G-Minor-Fontys/server/pkg/jwt"
	"github.com/google/uuid"
)

type PasswordHasher func(password string) (string, error)

func mapRegisterRequestToParams(hasher PasswordHasher, r *RegisterRequest) (*repository.RegisterParams, error) {
	passHash, err := hasher(r.Password)
	if err != nil {
		return nil, err
	}

	return &repository.RegisterParams{
		ID:           uuid.New(),
		Username:     r.Username,
		PasswordHash: passHash,
	}, nil
}

func mapResultToResponse(a *AuthenticationResult) *AuthenticationResponse {
	return &AuthenticationResponse{
		ID:          a.ID,
		AccessToken: a.AccessToken,
	}
}

func mapDatabaseRefreshTokenToToken(rt *repository.RefreshToken) *jwt.RefreshToken {
	return &jwt.RefreshToken{
		Token:      rt.Token,
		IdentityID: rt.IdentityID,
		CreatedAt:  rt.CreatedAt,
		ExpiresAt:  rt.ExpiresAt,
	}
}
