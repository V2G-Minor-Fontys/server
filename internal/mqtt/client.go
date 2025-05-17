package mqtt

import (
	"fmt"
	"github.com/V2G-Minor-Fontys/server/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

func NewClient(cfg *config.Mqtt) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port))
	opts.SetClientID("v2g-server")

	opts.OnConnect = func(client mqtt.Client) {
		slog.Info("Connected to MQTT server")
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
