package jwt

import (
	"crypto/rand"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AccessToken struct {
	Value     string    `json:"value,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func GenerateAccessToken(userID string, cfg *config.Jwt) (*AccessToken, error) {
	now := time.Now().UTC()
	exp := now.Add(cfg.Expire * time.Minute)

	claims := jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		Subject:   userID,
		Audience:  []string{cfg.Audience},
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(now),
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

func VerifyAccessToken(tokenStr string, cfg *config.Jwt) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	},
		jwt.WithIssuer(cfg.Issuer),
		jwt.WithAudience(cfg.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
		jwt.WithLeeway(time.Second*5),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims, nil
}

func GenerateRefreshToken() ([]byte, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
