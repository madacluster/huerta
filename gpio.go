package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type GPIOSensor struct {
	*Sensor
}

func NewGPIOSensor() *GPIOSensor {
	sensor := &GPIOSensor{NewSensor("Humidity", "20017", "20018", "20019", "20020")}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *GPIOSensor) getData() (float32, float32) {
	var humidity = core.Humidity{}
	messages := make(chan string)
	go s.keepAlivePort(messages)
	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &humidity)
	// Print Data
	fmt.Printf("Humidity: %f \t Temperature: %f\n", humidity.Humidity, humidity.Temperature)
	messages <- "Data Update Port: CONNECTED"
	// Start Base Port
	// go basePort() // Send Configuration Message

	// Close Data Update Port
	return humidity.Humidity, humidity.Temperature

}
