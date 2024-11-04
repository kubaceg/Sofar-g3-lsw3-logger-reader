package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/export/otlp"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/export/mosquitto"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug    bool `default:"false" yaml:"debug"`
	Inverter struct {
		Port          string   `yaml:"port"`
		LoggerSerial  uint     `yaml:"loggerSerial"`
		ReadInterval  int      `default:"60" yaml:"readInterval"`
		LoopLogging   bool     `default:"true" yaml:"loopLogging"`
		AttrWhiteList []string `yaml:"attrWhiteList"`
		AttrBlackList []string `yaml:"attrBlackList"`
	} `yaml:"inverter"`
	Mqtt mosquitto.MqttConfig `yaml:"mqtt"`
	Otlp otlp.Config          `yaml:"otlp"`
}

func (c *Config) validate() error {
	if c.Inverter.Port == "" {
		return errors.New("missing required inverter.port config")
	}

	if c.Inverter.LoggerSerial == 0 {
		return errors.New("missing required inverter.loggerSerial config")
	}

	return nil
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c Config) getLoglevel() slog.Leveler {
	if c.Debug {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}
