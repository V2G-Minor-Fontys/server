package charging_preferences

import (
	"encoding/json"
	"net/http"

	"github.com/V2G-Minor-Fontys/server/internal/config"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	"github.com/V2G-Minor-Fontys/server/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	svc *Service
}

func NewHandler(cfg *config.Jwt, db *pgxpool.Pool, queries *repository.Queries) *Handler {
	return &Handler{
		svc: NewService(cfg, db, queries),
	}
}

func (h *Handler) GetChargingPreferencesOfUserHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	idParam := chi.URLParam(r, "user-id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return httpx.BadRequest(ctx, "Invalid id")
	}

	preferences, err := h.svc.queries.ListChargingPreferencesForUser(r.Context(), id)
	if err != nil {
		return httpx.BadRequest(ctx, "Could not find preferences")
	}

	return json.NewEncoder(w).Encode(preferences)
}

func (h *Handler) CreateChargingPreferenceHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req ChargingPreference
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return httpx.BadRequest(ctx, "Could not parse JSON body")
	}

	a, _ := json.MarshalIndent(req, "", "  ")
	println(string(a))

	preferenceParams := repository.CreateChargingPreferenceParams{
		ID:       uuid.New(),
		UserID:   req.UserId,
		Name:     req.Name,
		Priority: req.Priority,
		Enabled:  req.Enabled,
	}

	if preferenceParams.Priority < 0 {
		return httpx.BadRequest(ctx, "Priority must be >= 0")
	}

	if req.OneTimeOccurrence != nil {
		if req.RegularOccurrence != nil {
			return httpx.BadRequest(ctx, "Either one_time_occurrence or regular_occurrence must be specified, but not both")
		}

		occurrenceParams, err := ToOneTimeOccurrenceParams(req.OneTimeOccurrence)
		if err != nil {
			return httpx.BadRequest(ctx, err.Error())
		}
		preferenceParams.OneTimeOccurrenceID.Scan(occurrenceParams.ID.String())

		err = h.svc.queries.CreateOneTimeOccurrence(ctx, occurrenceParams)
		if err != nil {
			println(err.Error())
			return httpx.BadRequest(ctx, "Could not create occurrence")
		}
	} else if req.RegularOccurrence != nil {

		occurrenceParams, err := ToRegularOccurrenceParams(req.RegularOccurrence)
		if err != nil {
			return httpx.BadRequest(ctx, err.Error())
		}
		preferenceParams.RegularOccurrenceID.Scan(occurrenceParams.ID.String())

		err = h.svc.queries.CreateRegularOccurrence(ctx, occurrenceParams)
		if err != nil {
			return httpx.BadRequest(ctx, "Could not create occurrence")
		}
	} else {
		return httpx.BadRequest(ctx, "Either one_time_occurrence or regular_occurrence must be specified")
	}

	if req.KeepChargeAt != nil {
		if req.ChargingPolicy != nil {
			return httpx.BadRequest(ctx, "automatic_charging must be specified if battery_charge is specified")
		}

		preferenceParams.KeepBatteryAt.Scan(int64(*req.KeepChargeAt))
	} else if req.ChargingPolicy != nil {
		if req.ChargingPolicy.MaxCharge < req.ChargingPolicy.MinCharge {
			return httpx.BadRequest(ctx, "max_charge must be greater than min_charge")
		}

		if req.ChargingPolicy.DischargeIfPriceAbove < req.ChargingPolicy.ChargeIfPriceBelow {
			return httpx.BadRequest(ctx, "discharge_if_price_above must be greater than charge_if_price_below")
		}

		chargeIfPriceBelow, err := ParseFloat(req.ChargingPolicy.ChargeIfPriceBelow)
		if err != nil {
			return httpx.BadRequest(ctx, err.Error())
		}
		dischargeIfPriceAbove, err := ParseFloat(req.ChargingPolicy.DischargeIfPriceAbove)
		if err != nil {
			return httpx.BadRequest(ctx, err.Error())
		}

		id := uuid.New()
		err = h.svc.queries.CreateChargingPolicies(ctx, repository.CreateChargingPoliciesParams{
			ID:                    id,
			MinCharge:             ParseInt(req.ChargingPolicy.MinCharge),
			MaxCharge:             ParseInt(req.ChargingPolicy.MaxCharge),
			ChargeIfPriceBelow:    chargeIfPriceBelow,
			DischargeIfPriceAbove: dischargeIfPriceAbove,
		})
		if err != nil {
			return httpx.BadRequest(ctx, "Could not create charging policies")
		}

		preferenceParams.ChargingPolicyID.Scan(id)
	} else {
		return httpx.BadRequest(ctx, "Either keep_charge_at or automatic_charging must be specified")
	}

	err := h.svc.queries.CreateChargingPreference(ctx, preferenceParams)
	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}

func (h *Handler) DeleteChargingPreferenceHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return httpx.BadRequest(ctx, "Invalid id")
	}

	err = h.svc.queries.DeleteChargingPreference(r.Context(), id)
	if err != nil {
		return httpx.BadRequest(ctx, "Could not charging preference")
	}

	return nil
}

func (h *Handler) CreateChargingPreferencesSchemaHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	idParam := chi.URLParam(r, "user-id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return httpx.BadRequest(ctx, "Invalid id")
	}

	preferences, err := h.svc.queries.ListChargingPreferencesForUser(r.Context(), id)
	if err != nil {
		return httpx.BadRequest(ctx, "Could not find preferences")
	}

	return json.NewEncoder(w).Encode(preferences)
}
