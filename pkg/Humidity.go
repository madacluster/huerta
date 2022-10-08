package pkg

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type HumiditySensor struct {
	*Sensor
}

func (s *HumiditySensor) SetData(config core.HumidityParams) (float32, float32) {
	driver := core.DriverConfig{
		Humidity: &config,
	}
	s.basePort(driver)
	return s.GetData()
}

func (s *HumiditySensor) Calibrate(temp float32) (float32, float32) {
	params := core.HumidityParams{
		CurrentTemperature: temp,
	}
	s.SetData(params)
	return s.GetData()
}

func NewHumiditySensor(host string) *HumiditySensor {
	sensor := &HumiditySensor{NewSensor("Humidity", host, 20017)}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *HumiditySensor) GetData() (float32, float32) {
	var humidity = core.Humidity{}
	messages := make(chan string)
	go s.keepAlivePort(messages)
	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &humidity)
	// Print Data
	log.Debug("Humidity: %f \t Temperature: %f\n", humidity.Humidity, humidity.Temperature)
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

func (s *HumiditySensor) Forward(c mqtt.Client) {
	for {
		humidity, temperature := s.GetData()
		s.forwardTemp(temperature, c)
		s.forwardHumidity(humidity, c)
		time.Sleep(300 * time.Second)
	}
}
