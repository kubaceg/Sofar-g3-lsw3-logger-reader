package ports

// support MQTT Discovery
type DiscoveryField struct {
	Name   string
	Factor string
	Unit   string
}

type Device interface {
	Name() string
	GetDiscoveryFields() []DiscoveryField
	Query() (MeasurementMap, error)
}

type MeasurementMap map[string]interface{}
