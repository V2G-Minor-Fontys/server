package controller

import (
	"github.com/google/uuid"
	"time"
)

type RegisterControllerRequest struct {
	CpuID           string `json:"cpuId"`
	FirmwareVersion string `json:"firmwareVersion"`
}

type PairUserToControllerRequest struct {
	UserId uuid.UUID `json:"userId"`
	CpuID  string    `json:"cpuId"`
}

type AddControllerTelemetryRequest struct {
	ControllerID  uuid.UUID `json:"controllerId"`
	TimeStamp     time.Time `json:"timeStamp"`
	OutputPower   int32     `json:"outputPower"`
	Soc           int16     `json:"soc"`
	EvDischarging bool      `json:"evDischarging"`
}

type UpdateControllerSettingsRequest struct {
	ID        uuid.UUID `json:"id"`
	Heartbeat int16     `json:"heartbeat"`
	AutoStart bool      `json:"autoStart"`
}

type Response struct {
	ID              uuid.UUID         `json:"id"`
	CpuID           string            `json:"cpuId"`
	FirmwareVersion string            `json:"firmwareVersion"`
	Settings        *SettingsResponse `json:"settings"`
}

type SettingsResponse struct {
	AutoStart     bool  `json:"autoStart"`
	HeartbeatRate int16 `json:"heartbeatRate"`
}

type TelemetryResponse struct {
	ControllerID  uuid.UUID `json:"controllerId"`
	Timestamp     time.Time `json:"timestamp"`
	OutputPower   int32     `json:"outputPower"`
	Soc           int16     `json:"soc"`
	EvDischarging bool      `json:"evDischarging"`
}
