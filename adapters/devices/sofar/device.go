package sofar

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

type Logger struct {
	serialNumber  uint
	connPort      ports.CommunicationPort
	attrWhiteList map[string]struct{}
	attrBlackList []*regexp.Regexp
}

// for a set in go we use a map of keys -> empty struct
func toSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{}, len(slice))
	v := struct{}{}
	for _, s := range slice {
		set[s] = v
	}
	return set
}

func toREs(patterns []string) []*regexp.Regexp {
	res := make([]*regexp.Regexp, 0, len(patterns))
	for idx, p := range patterns {
		re, err := regexp.Compile(p)
		if err == nil {
			res = append(res, re)
		} else {
			slog.Warn(fmt.Sprintf("config attrBlackList item %d '%s' not a valid regexp; %v", idx, p, err))
		}
	}
	return res
}

func NewSofarLogger(serialNumber uint, connPort ports.CommunicationPort, attrWhiteList []string, attrBlackList []string) *Logger {
	return &Logger{
		serialNumber:  serialNumber,
		connPort:      connPort,
		attrWhiteList: toSet(attrWhiteList),
		attrBlackList: toREs(attrBlackList),
	}
}

func (s *Logger) nameFilter(k string) bool {
	if len(s.attrWhiteList) > 0 {
		if _, ok := s.attrWhiteList[k]; ok {
			return true
		}
	}
	for _, re := range s.attrBlackList {
		if re.MatchString(k) {
			return false
		}
	}
	return true
}

func (s *Logger) GetDiscoveryFields() []ports.DiscoveryField {
	return getDiscoveryFields(s.nameFilter)
}

func (s *Logger) Query() (map[string]interface{}, error) {
	return readData(s.connPort, s.serialNumber, s.nameFilter)
}

func (s *Logger) Name() string {
	return "sofar"
}
