package filters

import "errors"

var ErrDailyGenerationDiffTooHigh = errors.New("daily generation spike detected, skipping")

type DailyGenerationSpikesFilter struct {
	lastDailyGenerationValue int
	maxDailyGenerationDiff   int
}

func NewDailyGenerationFilter(maxDailyGenerationDiff int) *DailyGenerationSpikesFilter {
	return &DailyGenerationSpikesFilter{
		maxDailyGenerationDiff: maxDailyGenerationDiff,
	}
}

func (d *DailyGenerationSpikesFilter) Filter(data map[string]interface{}) (map[string]interface{}, error) {
	if d.lastDailyGenerationValue > 0 &&
		data["PV_Generation_Today"].(int)-d.lastDailyGenerationValue > d.maxDailyGenerationDiff {
		return nil, ErrDailyGenerationDiffTooHigh
	}

	d.lastDailyGenerationValue = data["PV_Generation_Today"].(int)

	return data, nil
}
