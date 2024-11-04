package mosquitto

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

type MqttConfig struct {
	Url       string  `yaml:"url"`
	User      string  `yaml:"user"`
	Password  string  `yaml:"password"`
	Discovery *string `yaml:"ha_discovery_prefix,omitempty"`
	Prefix    string  `yaml:"prefix"`
}

type Connection struct {
	client mqtt.Client
	prefix string
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	slog.Debug("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	slog.Warn(fmt.Sprintf("Connect lost: %v", err))
}

func New(config *MqttConfig) (*Connection, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Url)
	opts.SetClientID("sofar")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	if config.User != "" {
		opts.SetUsername(config.User)
	}

	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	conn := &Connection{}
	conn.client = mqtt.NewClient(opts)
	conn.prefix = config.Prefix
	if token := conn.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return conn, nil

}

func (conn *Connection) publish(topic string, msg string, retain bool) {
	token := conn.client.Publish(topic, 0, retain, msg)
	res := token.WaitTimeout(1 * time.Second)
	if !res || token.Error() != nil {
		slog.Error(fmt.Sprintf("error inserting to MQTT: %s", token.Error()))
	}
}

// return "power" for kW etc., "energy" for kWh etc.
func unit2DeviceClass(unit string) string {
	if strings.HasSuffix(unit, "Wh") {
		return "energy"
	} else if strings.HasSuffix(unit, "W") {
		return "power"
	} else if strings.HasSuffix(unit, "Hz") {
		return "frequency"
	} else if strings.HasSuffix(unit, "VA") {
		return "apparent_power"
	} else if strings.HasSuffix(unit, "VAR") {
		return "reactive_power"
	} else if strings.HasSuffix(unit, "V") {
		return "voltage"
	} else if strings.HasSuffix(unit, "A") {
		return "current"
	} else if strings.HasSuffix(unit, "Ω") {
		return "voltage" // resistance not valid in https://developers.home-assistant.io/docs/core/entity/sensor/ so use "voltage"
	} else if strings.HasSuffix(unit, "℃") {
		return "temperature"
	} else if strings.HasSuffix(unit, "min") {
		return "duration"
	} else {
		return ""
	}
}

func unit2StateClass(unit string) string {
	if strings.HasSuffix(unit, "Wh") {
		return "total"
	} else {
		return "measurement"
	}
}

// MQTT Discovery: https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery
func (conn *Connection) InsertDiscoveryRecord(discovery string, state string, fields []ports.DiscoveryField) error {
	uniq := "01ad"
	for _, f := range fields {
		topic := fmt.Sprintf("%s/%s/config", discovery, f.Name)
		json, _ := json.Marshal(map[string]interface{}{
			"name":                  f.Name,
			"unique_id":             fmt.Sprintf("%s_%s", f.Name, uniq),
			"device_class":          unit2DeviceClass(f.Unit),
			"state_class":           unit2StateClass(f.Unit),
			"state_topic":           state,
			"unit_of_measurement":   f.Unit,
			"value_template":        fmt.Sprintf("{{ value_json.%s|int * %s }}", f.Name, f.Factor),
			"availability_topic":    state,
			"availability_template": "{{ value_json.availability }}",
			"device": map[string]interface{}{
				"identifiers": [...]string{fmt.Sprintf("Inverter_%s", uniq)},
				"name":        "Inverter",
			},
		})
		conn.publish(topic, string(json), true) // MQTT Discovery messages should be retained, but in dev it can become a pain
	}
	return nil
}

func (conn *Connection) InsertRecord(m map[string]interface{}) error {
	json, _ := json.Marshal(m)
	conn.publish(conn.prefix, string(json), false) // state messages should not be retained
	return nil
}

func (conn *Connection) Subscribe(topic string, callback mqtt.MessageHandler) {
	conn.client.Subscribe(topic, 0, callback)
}
