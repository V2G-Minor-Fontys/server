package controller

import (
	"github.com/google/uuid"
	"time"
)

type Controller struct {
	ID              uuid.UUID
	CpuID           string
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
	ID                  uuid.UUID
	ControllerID        uuid.UUID
	BatteryVoltage      float64
	BatteryCurrent      float64
	BatteryPower        float64
	BatteryState        int16
	InternalTemperature float64
	ModuleTemperature   float64
	RadiatorTemperature float64
	GridPowerR          int32
	TotalInverterPower  int32
	AcActivePower       int32
	LoadPowerR          int32
	TotalLoadPower      int32
	TotalEnergyToGrid   float64
	DailyEnergyToGrid   float64
	TotalEnergyFromGrid float64
	DailyEnergyFromGrid float64
	WorkMode            int16
	OperationMode       int16
	ErrorMessage        int64
	WarningCode         int16
}
