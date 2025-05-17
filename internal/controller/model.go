package controller

import (
	"github.com/google/uuid"
	"time"
)

type Controller struct {
	ID              uuid.UUID
	SerialNumber    string
	FirmwareVersion string
	Settings        *Settings
}

type Settings struct {
	ID            uuid.UUID
	AutoStart     bool
	HeartbeatRate int16
	UpdatedAt     time.Time
}

type Telemetry struct {
	ID            uuid.UUID
	ControllerID  uuid.UUID
	Timestamp     time.Time
	OutputPower   int32
	Soc           int16
	EvDischarging bool
}
