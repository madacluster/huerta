package pkg

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type UvSensor struct {
	*Sensor
}

func NewUvSensor(host string) *UvSensor {
	sensor := &UvSensor{NewSensor("UV", host, 20029)}
	return sensor
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *UvSensor) GetData() (float32, string) {
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

func (s *UvSensor) forwardIndex(temp float32, c mqtt.Client) {
	token := c.Publish("home/salon/uv_index", 0, false, fmt.Sprintf("%f", temp))
	token.Wait()
}
func (s *UvSensor) forwardRisk(temp string, c mqtt.Client) {
	token := c.Publish("home/salon/risk", 0, false, temp)
	token.Wait()
	if err := token.Error(); err != nil {
		fmt.Println(err)
	}
}

func (s *UvSensor) Forward(c mqtt.Client) {
	for {
		index, risk := s.GetData()
		s.forwardIndex(index, c)
		s.forwardRisk(risk, c)
		time.Sleep(300 * time.Second)
	}
}
