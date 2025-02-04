package filters

import (
	"testing"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
	"github.com/stretchr/testify/assert"
)

func TestFilterCases(t *testing.T) {
	tests := []struct {
		name                     string
		lastDailyGenerationValue uint32
		data                     ports.MeasurementMap
		expectedResult           ports.MeasurementMap
		expectedError            error
	}{
		{
			name:                     "ValidData",
			lastDailyGenerationValue: uint32(0),
			data:                     ports.MeasurementMap{"PV_Generation_Today": uint32(5000)},
			expectedResult:           ports.MeasurementMap{"PV_Generation_Today": uint32(5000)},
			expectedError:            nil,
		},
		{
			name:                     "DiffTooHigh",
			lastDailyGenerationValue: uint32(9000),
			data:                     ports.MeasurementMap{"PV_Generation_Today": uint32(20000)},
			expectedResult:           nil,
			expectedError:            ErrDailyGenerationDiffTooHigh,
		},
		{
			name:                     "FirstDataPoint",
			lastDailyGenerationValue: uint32(0),
			data:                     ports.MeasurementMap{"PV_Generation_Today": uint32(15000)},
			expectedResult:           ports.MeasurementMap{"PV_Generation_Today": uint32(15000)},
			expectedError:            nil,
		},
		{
			name:                     "ExactMaxDiff",
			lastDailyGenerationValue: uint32(10000),
			data:                     ports.MeasurementMap{"PV_Generation_Today": uint32(20000)},
			expectedResult:           ports.MeasurementMap{"PV_Generation_Today": uint32(20000)},
			expectedError:            nil,
		},
		{
			name:                     "NewDay",
			lastDailyGenerationValue: uint32(10000),
			data:                     ports.MeasurementMap{"PV_Generation_Today": uint32(500)},
			expectedResult:           ports.MeasurementMap{"PV_Generation_Today": uint32(500)},
			expectedError:            nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewDailyGenerationFilter(10000)
			filter.lastDailyGenerationValue = tt.lastDailyGenerationValue

			result, err := filter.Filter(tt.data)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
