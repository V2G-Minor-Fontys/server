package controller

import (
	"context"
	"github.com/V2G-Minor-Fontys/server/internal/mqtt"
)

func (h *Handler) MountMqttMessageHandlers() {
	h.mqtt.RegisterControllerTelemetrySubscriber(h.addControllerTelemetryMessageHandle)
}

func (h *Handler) addControllerTelemetryMessageHandle(msg *mqtt.AddControllerTelemetryMessage) error {
	return h.svc.AddControllerTelemetry(context.Background(), &AddControllerTelemetryRequest{
		ControllerID:        msg.ControllerID,
		BatteryVoltage:      msg.BatteryVoltage,
		BatteryCurrent:      msg.BatteryCurrent,
		BatteryPower:        msg.BatteryPower,
		BatteryState:        msg.BatteryState,
		InternalTemperature: msg.InternalTemperature,
		ModuleTemperature:   msg.ModuleTemperature,
		RadiatorTemperature: msg.RadiatorTemperature,
		GridPowerR:          msg.GridPowerR,
		TotalInverterPower:  msg.TotalInverterPower,
		AcActivePower:       msg.AcActivePower,
		LoadPowerR:          msg.LoadPowerR,
		TotalLoadPower:      msg.TotalLoadPower,
		TotalEnergyToGrid:   msg.TotalEnergyToGrid,
		DailyEnergyToGrid:   msg.DailyEnergyToGrid,
		TotalEnergyFromGrid: msg.TotalEnergyFromGrid,
		DailyEnergyFromGrid: msg.DailyEnergyFromGrid,
		WorkMode:            msg.WorkMode,
		OperationMode:       msg.OperationMode,
		ErrorMessage:        msg.ErrorMessage,
		WarningCode:         msg.WarningCode,
	})
}
