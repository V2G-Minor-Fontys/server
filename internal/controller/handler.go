package controller

import (
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/mqtt"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type Handler struct {
	mqtt *mqtt.Service
	svc  Service
}

func NewHandler(mqttService *mqtt.Service, db *pgxpool.Pool, queries *repository.Queries) *Handler {
	h := &Handler{
		mqtt: mqttService,
		svc:  NewService(db, queries),
	}
	return h
}

func (h *Handler) RegisterControllerHandler(w http.ResponseWriter, r *http.Request) error {
	var req RegisterControllerRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	controller, err := h.svc.RegisterController(r.Context(), &req)
	if err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusCreated, mapControllerToResponse(controller))
	return nil
}

func (h *Handler) PairUserToControllerHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := httpx.ParseUUIDParam(r, "userId")
	if err != nil {
		return err
	}

	var req PairUserToControllerRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	req.UserId = userID
	if err := h.svc.PairUserToController(r.Context(), &req); err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)
	return nil
}

func (h *Handler) GetControllerByCpuIdHandler(w http.ResponseWriter, r *http.Request) error {
	controller, err := h.svc.GetControllerByCpuId(r.Context(), chi.URLParam(r, "cpuId"))
	if err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusOK, mapControllerToResponse(controller))
	return nil
}

func (h *Handler) GetControllerTelemetryById(w http.ResponseWriter, r *http.Request) error {
	controllerID, err := httpx.ParseUUIDParam(r, "controllerId")
	if err != nil {
		return err
	}

	telemetry, err := h.svc.GetControllerTelemetryById(r.Context(), controllerID)
	if err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusOK, mapTelemetrySliceToResponse(telemetry))
	return nil
}

func (h *Handler) GetUserControllerHandler(w http.ResponseWriter, r *http.Request) error {
	userID, err := httpx.ParseUUIDParam(r, "userId")
	if err != nil {
		return err
	}

	controller, err := h.svc.GetControllerByUserId(r.Context(), userID)
	if err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusOK, mapControllerToResponse(controller))
	return nil
}

func (h *Handler) UpdateControllerSettingsHandler(w http.ResponseWriter, r *http.Request) error {
	controllerID, err := httpx.ParseUUIDParam(r, "controllerId")
	if err != nil {
		return err
	}

	var req UpdateControllerSettingsRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	req.ID = controllerID
	if err := h.svc.UpdateControllerSettings(r.Context(), &req); err != nil {
		return err
	}

	if err := h.mqtt.UpdateControllerSettings(r.Context(), mqtt.UpdateControllerSettings{
		ControllerID: req.ID,
		Heartbeat:    req.Heartbeat,
		AutoStart:    req.AutoStart,
	}); err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)
	return nil
}

func (h *Handler) ExecuteControllerActionHandler(w http.ResponseWriter, r *http.Request) error {
	controllerId, err := httpx.ParseUUIDParam(r, "controllerId")
	if err != nil {
		return err
	}

	var req mqtt.ControllerActionRequest
	if err := httpx.DecodeJSONBody(r, &req); err != nil {
		return err
	}

	req.ControllerID = controllerId
	if err := h.mqtt.ExecuteControllerAction(r.Context(), &req); err != nil {
		return err
	}

	httpx.ResponseWithJSON(w, http.StatusNoContent, nil)
	return nil
}

func (h *Handler) ShutdownMQTT(timeoutMils uint) {
	h.mqtt.ShutdownMQTT(timeoutMils)
}
