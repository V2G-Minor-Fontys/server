package controller

import (
	"github.com/V2G-Minor-Fontys/server/internal/repository"
)

func mapDatabaseTelemetrySliceToTelemetrySlice(data []repository.ControllerTelemetry) []*Telemetry {
	result := make([]*Telemetry, 0, len(data))
	for _, item := range data {
		result = append(result, mapDatabaseTelemetryToTelemetry(&item))
	}
	return result
}

func mapDatabaseTelemetryToTelemetry(telemetry *repository.ControllerTelemetry) *Telemetry {
	controllerID, _ := repository.ParsePgUUIDToGuid(telemetry.ControllerID)
	return &Telemetry{
		ID:            telemetry.ID,
		ControllerID:  controllerID,
		Timestamp:     telemetry.Timestamp,
		OutputPower:   telemetry.OutputPower,
		Soc:           telemetry.Soc,
		EvDischarging: telemetry.EvDischarging,
	}
}

func mapSettingsToResponse(settings *Settings) *SettingsResponse {
	return &SettingsResponse{
		AutoStart:     settings.AutoStart,
		HeartbeatRate: settings.HeartbeatRate,
	}
}

func mapControllerToResponse(c *Controller) *Response {
	return &Response{
		ID:              c.ID,
		SerialNumber:    c.SerialNumber,
		FirmwareVersion: c.FirmwareVersion,
		Settings:        mapSettingsToResponse(c.Settings),
	}
}

func mapDatabaseSettingsToSettings(s *repository.ControllerSetting) *Settings {
	return &Settings{
		ID:            s.ID,
		AutoStart:     s.AutoStart,
		HeartbeatRate: s.HeartbeatRate,
		UpdatedAt:     s.UpdatedAt,
	}
}

func mapTelemetrySliceToResponse(data []*Telemetry) []*TelemetryResponse {
	result := make([]*TelemetryResponse, 0, len(data))
	for _, item := range data {
		result = append(result, mapTelemetryToResponse(item))
	}
	return result
}

func mapTelemetryToResponse(telemetry *Telemetry) *TelemetryResponse {
	return &TelemetryResponse{
		ControllerID:  telemetry.ControllerID,
		Timestamp:     telemetry.Timestamp,
		OutputPower:   telemetry.OutputPower,
		Soc:           telemetry.Soc,
		EvDischarging: telemetry.EvDischarging,
	}
}
