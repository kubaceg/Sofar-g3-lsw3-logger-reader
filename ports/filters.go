package ports

type Filter interface {
	Filter(map[string]interface{}) (map[string]interface{}, error)
}
