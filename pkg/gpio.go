package pkg

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

type GPIOSensor struct {
	*Sensor
}

func NewGPIOSensor(host string) *GPIOSensor {
	sensor := &GPIOSensor{NewSensor("GPRIO", host, 20049)}
	return sensor
}

func (s *GPIOSensor) SetData(config core.GpioParams) ([]int, []int) {
	driver := core.DriverConfig{
		Gpio: &config,
	}
	s.basePort(driver)
	return s.GetData()
}

func (s *GPIOSensor) ReadPin(pin int) bool {
	on, _ := s.GetData()
	for _, v := range on {
		if v == pin {
			return true
		}
	}
	return false
}

func (s *GPIOSensor) WritePin(pin int, value int) ([]int, []int) {
	var gpio = core.GpioParams{}
	gpio.Pin = uint32(pin)
	gpio.Mode = core.GpioParams_OUTPUT
	gpio.Value = uint32(value)
	return s.SetData(gpio)
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *GPIOSensor) GetData() ([]int, []int) {
	var gpio = core.GpioParams{}
	messages := make(chan string)
	go s.keepAlivePort(messages)
	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &gpio)
	// Print Data
	fmt.Println(gpio)
	fmt.Printf("Pin: %d\nMode: %v\nValue: %d\nValues: %d\n", gpio.Pin, gpio.Mode, gpio.Value, gpio.Values)
	messages <- "Data Update Port: CONNECTED"

	// Start Base Port
	// go basePort() // Send Configuration Message
	listValues := Reverse(fmt.Sprintf("%028s", strconv.FormatInt(int64(gpio.Values), 2)))
	// Close Data Update Port
	on := make([]int, 0)
	off := make([]int, 0)
	for i, v := range listValues {
		if v == '1' {
			on = append(on, i)
		} else {
			off = append(off, i)
		}
	}
	return on, off

}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
