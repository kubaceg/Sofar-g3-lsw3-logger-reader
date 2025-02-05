package ports

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Database interface {
	InsertDiscoveryRecord(discovery string, prefix string, fields []DiscoveryField) error
	InsertRecord(measurement MeasurementMap) error
}

type DatabaseWithListener interface {
	Database
	Subscribe(topic string, callback mqtt.MessageHandler)
}
