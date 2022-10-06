package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type PressureSensor struct {
	*Sensor
}

func NewPressureSensor() *PressureSensor {
	sensor := &PressureSensor{NewSensor("Pressure", "20025", "20026", "20027", "20028")}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *PressureSensor) getData() (float32, float32, float32) {
	var pressure = core.Pressure{}
	messages := make(chan string)
	go s.keepAlivePort(messages)

	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &pressure)
	// Print Data
	fmt.Printf("pressure: %f\taltitude: %f\ttemperature: %f\n", pressure.Pressure, pressure.Altitude, pressure.Temperature)
	messages <- "Data Update Port: CONNECTED"
	// Start Base Port
	// go basePort() // Send Configuration Message

	// Close Data Update Port
	return pressure.Pressure, pressure.Altitude, pressure.Temperature

}
