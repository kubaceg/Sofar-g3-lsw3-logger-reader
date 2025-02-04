package ports

type Filter interface {
	Filter(MeasurementMap) (MeasurementMap, error)
}
