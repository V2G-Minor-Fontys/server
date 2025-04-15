package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

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
