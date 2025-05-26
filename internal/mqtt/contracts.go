package mqtt

import (
	"github.com/google/uuid"
)

type AddControllerTelemetryMessage struct {
	ControllerID        uuid.UUID `json:"controllerId"`
	BatteryVoltage      float64   `json:"batteryVoltage"`
	BatteryCurrent      float64   `json:"batteryCurrent"`
	BatteryPower        float64   `json:"batteryPower"`
	BatteryState        int16     `json:"batteryState"`
	InternalTemperature float64   `json:"internalTemperature"`
	ModuleTemperature   float64   `json:"moduleTemperature"`
	RadiatorTemperature float64   `json:"radiatorTemperature"`
	GridPowerR          int32     `json:"gridPowerR"`
	TotalInverterPower  int32     `json:"totalInverterPower"`
	AcActivePower       int32     `json:"acActivePower"`
	LoadPowerR          int32     `json:"loadPowerR"`
	TotalLoadPower      int32     `json:"totalLoadPower"`
	TotalEnergyToGrid   float64   `json:"totalEnergyToGrid"`
	DailyEnergyToGrid   float64   `json:"dailyEnergyToGrid"`
	TotalEnergyFromGrid float64   `json:"totalEnergyFromGrid"`
	DailyEnergyFromGrid float64   `json:"dailyEnergyFromGrid"`
	WorkMode            int16     `json:"workMode"`
	OperationMode       int16     `json:"operationMode"`
	ErrorMessage        int64     `json:"errorMessage"`
	WarningCode         int16     `json:"warningCode"`
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
