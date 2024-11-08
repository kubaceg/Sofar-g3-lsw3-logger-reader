package filters

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCases(t *testing.T) {
	tests := []struct {
		name                     string
		lastDailyGenerationValue uint32
		data                     map[string]interface{}
		expectedResult           map[string]interface{}
		expectedError            error
	}{
		{
			name:                     "ValidData",
			lastDailyGenerationValue: 0,
			data:                     map[string]interface{}{"PV_Generation_Today": uint32(5000)},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": uint32(5000)},
			expectedError:            nil,
		},
		{
			name:                     "DiffTooHigh",
			lastDailyGenerationValue: 9000,
			data:                     map[string]interface{}{"PV_Generation_Today": uint32(20000)},
			expectedResult:           nil,
			expectedError:            errors.New("PV generation today diff is too high, skipping"),
		},
		{
			name:                     "FirstDataPoint",
			lastDailyGenerationValue: 0,
			data:                     map[string]interface{}{"PV_Generation_Today": uint32(15000)},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": uint32(15000)},
			expectedError:            nil,
		},
		{
			name:                     "ExactMaxDiff",
			lastDailyGenerationValue: 10000,
			data:                     map[string]interface{}{"PV_Generation_Today": uint32(20000)},
			expectedResult:           map[string]interface{}{"PV_Generation_Today": uint32(20000)},
			expectedError:            nil,
		},
		{
			name:                     "Nil",
			lastDailyGenerationValue: 10000,
			data:                     map[string]interface{}{"PV_Generation_Today": nil},
			expectedResult:           nil,
			expectedError:            errors.New("PV generation today not found, skipping"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewDailyGenerationFilter()
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
