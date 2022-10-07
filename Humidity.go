package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type HumiditySensor struct {
	*Sensor
}

func NewHumiditySensor() *HumiditySensor {
	sensor := &HumiditySensor{NewSensor("Humidity", "20017", "20018", "20019", "20020")}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *HumiditySensor) getData() (float32, float32) {
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

func (s *HumiditySensor) forwardTemp(temp float32, c mqtt.Client) {
	token := c.Publish("home/salon/temperature", 0, false, fmt.Sprintf("%f", temp))
	fmt.Println(token.Wait())
}
func (s *HumiditySensor) forwardHumidity(temp float32, c mqtt.Client) {
	token := c.Publish("home/salon/humidity", 0, false, fmt.Sprintf("%f", temp))

	fmt.Println(token.Wait())
	if err := token.Error(); err != nil {
		fmt.Println(err)
	}
}

func (s *HumiditySensor) forward(c mqtt.Client) {
	for {
		humidity, temperature := s.getData()
		s.forwardTemp(temperature, c)
		s.forwardHumidity(humidity, c)
		time.Sleep(300 * time.Second)
	}
}
