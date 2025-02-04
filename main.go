package main

import (
	"context"
	"fmt"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	gser "go.bug.st/serial"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/comms/serial"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/comms/tcpip"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/devices/sofar"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/export/mosquitto"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/export/otlp"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/filters"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/ports"
)

const maximumFailedConnections = 3

var (
	config      *Config
	port        ports.CommunicationPort
	mqtt        ports.DatabaseWithListener
	device      ports.Device
	telem       *otlp.Service
	dataFilters []ports.Filter

	hasMQTT bool
	hasOTLP bool
)

func initialize() {
	var err error
	config, err = NewConfig("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("error during config.yaml file load: %s", err))
	}

	initializeLogger()
	initializeCommunicationPorts()
	initializeServices()
	initializeDevice()
	initializeFilters()
}

func initializeLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: config.getLoglevel()}))
	slog.SetDefault(logger)
}

func initializeCommunicationPorts() {
	if isSerialPort(config.Inverter.Port) {
		port = serial.New(config.Inverter.Port, 2400, 8, gser.NoParity, gser.OneStopBit)
		slog.Debug(fmt.Sprintf("using serial communications port %s", config.Inverter.Port))
	} else {
		port = tcpip.New(config.Inverter.Port)
		slog.Debug(fmt.Sprintf("using TCP/IP communications port %s", config.Inverter.Port))
	}
}

func initializeServices() {
	var err error
	hasMQTT = config.Mqtt.Url != "" && config.Mqtt.Prefix != ""
	hasOTLP = config.Otlp.Grpc.Url != "" || config.Otlp.Http.Url != ""

	if hasMQTT {
		mqtt, err = mosquitto.New(&config.Mqtt)
		if err != nil {
			slog.Error(fmt.Sprintf("MQTT connection failed: %s", err))
			os.Exit(1)
		}
	}

	if hasOTLP {
		telem, err = otlp.New(&config.Otlp)
		if err != nil {
			slog.Error(fmt.Sprintf("error initializing otlp connection: %s", err))
			os.Exit(1)
		}
	}
}

func initializeDevice() {
	device = sofar.NewSofarLogger(config.Inverter.LoggerSerial, port, config.Inverter.AttrWhiteList, config.Inverter.AttrBlackList)
}

func initializeFilters() {
	if config.Filters.DailyGenerationSpikes > 0 {
		dataFilters = []ports.Filter{filters.NewDailyGenerationFilter(config.Filters.DailyGenerationSpikes)}
	}
}

func main() {
	initialize()

	if hasMQTT && config.Mqtt.Discovery == nil {
		_ = mqtt.InsertDiscoveryRecord(*config.Mqtt.Discovery, config.Mqtt.Prefix, device.GetDiscoveryFields()) // logs errors, always returns nil
	}

	for {
		performMeasurements()
		time.Sleep(time.Duration(config.Inverter.ReadInterval) * time.Second)
	}
}

func performMeasurements() {
	if config.Inverter.LoopLogging {
		slog.Debug("performing measurements")
	}

	var measurements map[string]interface{}
	var err error
	for retry := 0; measurements == nil && retry < maximumFailedConnections; retry++ {
		measurements, err = device.Query()
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to perform measurements on retry %d: %s", retry, err))
			// at night, inverter is offline, err = "dial tcp 192.168.xx.xxx:8899: i/o timeout"
			// at other times occasionally: "read tcp 192.168.68.104:38670->192.168.68.106:8899: i/o timeout"
			return
		}
	}

	for _, filter := range dataFilters {
		measurements, err = filter.Filter(measurements)
		if err != nil {
			slog.Error(fmt.Sprintf("filter failed: %v", err))
			return
		}
	}

	if hasMQTT {
		publishToMQTT(measurements)
	}

	if hasOTLP && measurements != nil {
		err := telem.CollectAndPushMetrics(context.Background(), measurements)
		if err != nil {
			slog.Error(fmt.Sprintf("error recording telemetry: %s", err))
		} else {
			slog.Debug("measurements pushed via OLTP")
		}
	}
}

func publishToMQTT(measurements map[string]interface{}) {
	var m map[string]interface{}
	timeStamp := time.Now().UnixNano() / int64(time.Millisecond)
	if measurements != nil {
		m = make(map[string]interface{}, len(measurements)+2)
		for k, v := range measurements {
			m[k] = v
		}
		m["availability"] = "online"
		m["LastTimestamp"] = timeStamp
	} else {
		m = map[string]interface{}{
			"availability":  "offline",
			"LastTimestamp": timeStamp,
		}
	}
	_ = mqtt.InsertRecord(m) // logs errors, always returns nil
}

func isSerialPort(portName string) bool {
	return strings.HasPrefix(portName, "/")
}
