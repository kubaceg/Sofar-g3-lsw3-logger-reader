package sofar

import (
	"testing"
	"regexp"
)

func TestNameFilter(t *testing.T) {
	// Create a Logger instance with the desired attributes
	logger := &Logger{
		attrWhiteList: map[string]struct{}{
			"whitelisted": {},
		},
		attrBlackList: []*regexp.Regexp{
			regexp.MustCompile("^blacklisted"),
		},
	}

	// Test case 1: Key in the white list
	result := logger.nameFilter("whitelisted")
	if result != true {
		t.Errorf("Expected: true, Got: %v", result)
	}

	// Test case 2: Key not in the white list, but not matching any black list regex
	result = logger.nameFilter("notblacklisted")
	if result != true {
		t.Errorf("Expected: true, Got: %v", result)
	}

	// Test case 3: Key in the black list
	result = logger.nameFilter("blacklisted-key")
	if result != false {
		t.Errorf("Expected: false, Got: %v", result)
	}

	// Test case 4: Key not in the white list and matches a black list regex
	result = logger.nameFilter("blacklisted")
	if result != false {
		t.Errorf("Expected: false, Got: %v", result)
	}

	// Test case 5: No white or black list
	logger = &Logger{} // Reset the logger
	result = logger.nameFilter("anykey")
	if result != true {
		t.Errorf("Expected: true, Got: %v", result)
	}
}
