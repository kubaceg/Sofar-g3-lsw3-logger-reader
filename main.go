package main

import (
	"context"
	"log"
	_ "net/http/pprof"
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

// maximumFailedConnections maximum number failed logger connection, after this number will be exceeded reconnect
// interval will be extended from 5s to readInterval defined in config file
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
		log.Fatalln(err)
	}

	hasMQTT = config.Mqtt.Url != "" && config.Mqtt.Prefix != ""
	hasOTLP = config.Otlp.Grpc.Url != "" || config.Otlp.Http.Url != ""

	if isSerialPort(config.Inverter.Port) {
		port = serial.New(config.Inverter.Port, 2400, 8, gser.NoParity, gser.OneStopBit)
		log.Printf("using serial communcations port %s", config.Inverter.Port)
	} else {
		port = tcpip.New(config.Inverter.Port)
		log.Printf("using TCP/IP communications port %s", config.Inverter.Port)
	}

	if hasMQTT {
		mqtt, err = mosquitto.New(&config.Mqtt)
		if err != nil {
			log.Fatalf("MQTT connection failed: %s", err)
		}

	}

	if hasOTLP {
		telem, err = otlp.New(&config.Otlp)
		if err != nil {
			log.Fatalf("error initializating otlp connection: %s", err)
		}
	}

	dataFilters = []ports.Filter{filters.NewDailyGenerationFilter()}

	device = sofar.NewSofarLogger(config.Inverter.LoggerSerial, port)
}

func main() {
	initialize()
	failedConnections := 0

	for {
		log.Printf("performing measurements")
		timeStart := time.Now()

		measurements, err := device.Query()
		if err != nil {
			log.Printf("failed to perform measurements: %s", err)
			failedConnections++

			if failedConnections > maximumFailedConnections {
				time.Sleep(time.Duration(config.Inverter.ReadInterval) * time.Second)
			}

			continue
		}

		for _, filter := range dataFilters {
			measurements, err = filter.Filter(measurements)
			if err != nil {
				log.Printf("filter failed: %s", err)
				continue
			}
		}

		failedConnections = 0

		if hasMQTT {
			go func() {
				err = mqtt.InsertRecord(measurements)
				if err != nil {
					log.Printf("failed to insert record to MQTT: %s\n", err)
				} else {
					log.Println("measurements pushed to MQTT")
				}
			}()
		}

		if hasOTLP {
			go func() {
				err = telem.CollectAndPushMetrics(context.Background(), measurements)
				if err != nil {
					log.Printf("error recording telemetry: %s\n", err)
				} else {
					log.Println("measurements pushed via OLTP")
				}
			}()
		}

		duration := time.Since(timeStart)

		delay := time.Duration(config.Inverter.ReadInterval)*time.Second - duration
		if delay <= 0 {
			delay = 1 * time.Second
		}

		time.Sleep(delay)
	}

}

func isSerialPort(portName string) bool {
	return strings.HasPrefix(portName, "/")
}
