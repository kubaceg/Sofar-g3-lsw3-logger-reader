package filters

import "errors"

const maxDailyGenerationDiff = 10000

type DailyGenerationFilter struct {
	lastDailyGenerationValue int
}

func NewDailyGenerationFilter() *DailyGenerationFilter {
	return &DailyGenerationFilter{}
}

func (d *DailyGenerationFilter) Filter(data map[string]interface{}) (map[string]interface{}, error) {
	if d.lastDailyGenerationValue > 0 &&
		data["PV_Generation_Today"].(int)-d.lastDailyGenerationValue > maxDailyGenerationDiff {
		return nil, errors.New("PV generation today diff is too high, skipping")
	}

	d.lastDailyGenerationValue = data["PV_Generation_Today"].(int)

	return data, nil
}
