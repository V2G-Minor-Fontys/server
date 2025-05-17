package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/V2G-Minor-Fontys/server/internal/httpx"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"strings"
	"time"
)

type Service struct {
	client   mqtt.Client
	Shutdown func(timeoutMils uint)
}

func NewService(client mqtt.Client) *Service {
	return &Service{
		client:   client,
		Shutdown: client.Disconnect,
	}
}

func (h *Service) RegisterControllerTelemetrySubscriber(onControllerTelemetryReceived func(msg *AddControllerTelemetryMessage) error) {
	h.client.Subscribe("v2g/controllers/telemetry", 0, func(client mqtt.Client, message mqtt.Message) {
		defer message.Ack()
		var msg AddControllerTelemetryMessage
		payload := message.Payload()

		if err := json.Unmarshal(payload, &msg); err != nil {
			slog.Error("Failed to unmarshal telemetry message",
				"controller_id", msg.ControllerID,
				"error", err,
				"topic", message.Topic(),
				"payload", string(payload),
			)
			return
		}

		if err := onControllerTelemetryReceived(&msg); err != nil {
			slog.Error("Failed to consume telemetry message",
				"controller_id", msg.ControllerID,
				"error", err,
				"topic", message.Topic(),
				"payload", string(payload),
			)
		}
	})
}

func (h *Service) ExecuteControllerAction(ctx context.Context, req ControllerActionRequest) error {
	switch strings.ToLower(req.Action) {
	case "start_discharging", "stop_discharging":
		payload, err := json.Marshal(req)
		if err != nil {
			return httpx.InternalErr(ctx, "Failed to encode request payload", err)
		}

		token := h.client.Publish("v2g/controller/action", 0, false, payload)
		if !token.WaitTimeout(250 * time.Millisecond) {
			return httpx.InternalErr(ctx, "Timeout occurred while publishing MQTT message", errors.New("MQTT server is not responding fast enough"))
		}

		if token.Error() != nil {
			return httpx.InternalErr(ctx, "Failed to publish MQTT message", token.Error())
		}

		return nil

	default:
		return httpx.BadRequest(ctx, fmt.Sprintf("Unsupported action: %s. Currently Supported: %s, %s", req.Action, "start_discharging", "stop_discharging"))
	}
}
