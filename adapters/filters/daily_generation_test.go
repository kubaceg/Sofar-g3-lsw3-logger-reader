package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCases(t *testing.T) {
	tests := []struct {
		name                     string
		lastDailyGenerationValue int
		data                     map[string]interface{}
		expectedResult           map[string]interface{}
		expectedError            error
	}{
		{
			name:                     "ValidData",
			lastDailyGenerationValue: 0,
			data:                     map[string]interface{}{"PV_Generation_Today": 5000},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": 5000},
			expectedError:            nil,
		},
		{
			name:                     "DiffTooHigh",
			lastDailyGenerationValue: 9000,
			data:                     map[string]interface{}{"PV_Generation_Today": 20000},
			expectedResult:           nil,
			expectedError:            ErrDailyGenerationDiffTooHigh,
		},
		{
			name:                     "FirstDataPoint",
			lastDailyGenerationValue: 0,
			data:                     map[string]interface{}{"PV_Generation_Today": 15000},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": 15000},
			expectedError:            nil,
		},
		{
			name:                     "ExactMaxDiff",
			lastDailyGenerationValue: 10000,
			data:                     map[string]interface{}{"PV_Generation_Today": 20000},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": 20000},
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
