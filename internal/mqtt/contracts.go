package mqtt

import (
	"github.com/google/uuid"
	"time"
)

type AddControllerTelemetryMessage struct {
	ControllerID  uuid.UUID `json:"controllerId"`
	TimeStamp     time.Time `json:"timeStamp"`
	OutputPower   int       `json:"outputPower"`
	Soc           int16     `json:"soc"`
	EvDischarging bool      `json:"evDischarging"`
}

type ControllerActionRequest struct {
	ControllerID uuid.UUID `json:"controllerId,omitempty"`
	Action       string    `json:"action"`
}

type UpdateControllerSettings struct {
	ControllerID uuid.UUID `json:"controllerId"`
	Heartbeat    int16     `json:"heartbeat"`
	AutoStart    bool      `json:"autoStart"`
}
