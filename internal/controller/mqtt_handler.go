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
		ControllerID:  msg.ControllerID,
		TimeStamp:     msg.TimeStamp,
		OutputPower:   int32(msg.OutputPower),
		Soc:           msg.Soc,
		EvDischarging: msg.EvDischarging,
	})
}
