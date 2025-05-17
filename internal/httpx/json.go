package httpx

import (
	"encoding/hex"
	"encoding/json"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"log/slog"
	"net/http"
)

func DecodeJSONBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return BadRequest(r.Context(), "Could not parse JSON body")
	}
	return nil
}

func ResponseWithJSON(w http.ResponseWriter, statusCode int, content any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if content != nil {
		jsonContent, _ := json.Marshal(content)
		_, err := w.Write(jsonContent)
		if err != nil {
			slog.Error("Error writing response body", "error", err)
		}
	}
}

func SetRefreshToken(w http.ResponseWriter, refreshToken *repository.RefreshToken) {
	cookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	if refreshToken != nil {
		cookie.MaxAge = 0
		cookie.Value = hex.EncodeToString(refreshToken.Token)
		cookie.Expires = refreshToken.ExpiresAt
	}

	http.SetCookie(w, cookie)
}
