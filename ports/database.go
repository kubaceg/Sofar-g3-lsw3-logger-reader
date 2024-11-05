package ports

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Database interface {
	InsertDiscoveryRecord(discovery string, prefix string, fields []DiscoveryField) error
	InsertRecord(measurement map[string]interface{}) error
}

type DatabaseWithListener interface {
	Database
	Subscribe(topic string, callback mqtt.MessageHandler)
}
