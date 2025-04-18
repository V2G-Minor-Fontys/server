package system

import (
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) error {
	httpx.ResponseWithJSON(w, http.StatusOK, "Healthy")
	return nil
}
