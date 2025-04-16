package jwt

import (
	"crypto/rand"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
}

type AccessToken struct {
	Value     string    `json:"value,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func GenerateAccessToken(userID string, cfg *config.Jwt) (*AccessToken, error) {
	now := time.Now().UTC()
	exp := now.Add(cfg.Expire * time.Minute)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   userID,
			Audience:  []string{cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Value:     tokenString,
		ExpiresAt: exp,
	}, nil
}

func GenerateRefreshToken() ([]byte, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
