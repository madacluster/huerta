package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type UvSensor struct {
	*Sensor
}

func NewUvSensor() *UvSensor {
	sensor := &UvSensor{NewSensor("Pressure", "20029", "20030", "20031", "20032")}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *UvSensor) getData() (float32, string) {
	var uv = core.UV{}
	messages := make(chan string)
	go s.keepAlivePort(messages)
	// time.Sleep(1000 * time.Millisecond)
	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &uv)
	// Print Data
	fmt.Printf("index: %f\trisk: %s\t\n", uv.UvIndex, uv.OmsRisk)
	messages <- "Data Update Port: CONNECTED"
	// Start Base Port
	// go basePort() // Send Configuration Message

	// Close Data Update Port
	return uv.UvIndex, uv.OmsRisk

}
