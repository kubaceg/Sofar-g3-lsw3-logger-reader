package filters

import (
	"errors"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

var ErrDailyGenerationDiffTooHigh = errors.New("daily generation spike detected, skipping")

type DailyGenerationSpikesFilter struct {
	lastDailyGenerationValue uint32
	maxDailyGenerationDiff   uint32
}

func NewDailyGenerationFilter(maxDailyGenerationDiff uint32) *DailyGenerationSpikesFilter {
	return &DailyGenerationSpikesFilter{
		maxDailyGenerationDiff: maxDailyGenerationDiff,
	}
}

func (d *DailyGenerationSpikesFilter) Filter(data ports.MeasurementMap) (ports.MeasurementMap, error) {
	// If the last daily generation value is higher than the current one, it means a new day has started
	if data["PV_Generation_Today"].(uint32) < d.lastDailyGenerationValue {
		d.lastDailyGenerationValue = 0
	}

	if d.lastDailyGenerationValue > 0 &&
		data["PV_Generation_Today"].(uint32)-d.lastDailyGenerationValue > d.maxDailyGenerationDiff {
		return nil, ErrDailyGenerationDiffTooHigh
	}

	d.lastDailyGenerationValue = data["PV_Generation_Today"].(uint32)

	return data, nil
}
