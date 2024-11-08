package filters

import "errors"

const maxDailyGenerationDiff = 10000

type DailyGenerationFilter struct {
	lastDailyGenerationValue uint32
}

func NewDailyGenerationFilter() *DailyGenerationFilter {
	return &DailyGenerationFilter{}
}

func (d *DailyGenerationFilter) Filter(data map[string]interface{}) (map[string]interface{}, error) {
	var todayGeneration uint32
	var ok bool
	if todayGeneration, ok = data["PV_Generation_Today"].(uint32); !ok {
		return nil, errors.New("PV generation today not found, skipping")
	}

	if d.lastDailyGenerationValue > 0 &&
		todayGeneration-d.lastDailyGenerationValue > maxDailyGenerationDiff {
		return nil, errors.New("PV generation today diff is too high, skipping")
	}

	d.lastDailyGenerationValue = todayGeneration

	return data, nil
}
