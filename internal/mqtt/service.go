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
	client mqtt.Client
}

func NewService(client mqtt.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) publishMQTTMessage(ctx context.Context, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return httpx.InternalErr(ctx, "Failed to encode request payload", err)
	}

	token := s.client.Publish(topic, 0, false, data)
	if !token.WaitTimeout(250 * time.Millisecond) {
		return httpx.InternalErr(ctx, "Timeout occurred while publishing MQTT message", errors.New("MQTT server is not responding fast enough"))
	}

	if token.Error() != nil {
		return httpx.InternalErr(ctx, "Failed to publish MQTT message", token.Error())
	}

	return nil
}

func (s *Service) RegisterControllerTelemetrySubscriber(onControllerTelemetryReceived func(msg *AddControllerTelemetryMessage) error) {
	s.client.Subscribe("v2g/controller/telemetry", 0, func(client mqtt.Client, message mqtt.Message) {
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

func (s *Service) ExecuteControllerAction(ctx context.Context, req *ControllerActionRequest) error {
	switch strings.ToLower(req.Action) {
	case "start_discharging", "stop_discharging":
		return s.publishMQTTMessage(ctx, "v2g/controller/action", req)
	default:
		return httpx.BadRequest(ctx, fmt.Sprintf("Unsupported action: %s. Currently Supported: %s, %s", req.Action, "start_discharging", "stop_discharging"))
	}
}

func (s *Service) UpdateControllerSettings(ctx context.Context, req UpdateControllerSettings) error {
	return s.publishMQTTMessage(ctx, "v2g/controller/settings", req)
}

func (s *Service) ShutdownMQTT(timeoutMils uint) {
	s.client.Disconnect(timeoutMils)
}
