package pkg

import (
	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
	log "github.com/sirupsen/logrus"
)

type EverloopSensor struct {
	*Sensor
}

type LED struct {
	White, Red, Blue, Green int
}

func NewEverloopSensor(host string) *EverloopSensor {
	sensor := &EverloopSensor{NewSensor("Everloop", host, 20021)}
	return sensor
}

func (s *EverloopSensor) SetData(config core.EverloopImage) (int, []*core.LedValue) {
	driver := core.DriverConfig{
		Image: &config,
	}
	s.basePort(driver)
	return s.GetData()
}

func (s *EverloopSensor) WriteLed(leds []*LED) {
	coreLeds := make([]*core.LedValue, 0)
	for _, l := range leds {
		led := &core.LedValue{
			Blue:  uint32(l.Blue),
			Red:   uint32(l.Red),
			Green: uint32(l.Green),
			White: uint32(l.White),
		}
		coreLeds = append(coreLeds, led)
	}
	config := core.EverloopImage{
		Led: coreLeds,
	}
	s.SetData(config)
}

// DATA UPDATE PORT \\ (port where updates are received)
func (s *EverloopSensor) GetData() (int, []*core.LedValue) {
	var everloop = core.EverloopImage{}
	messages := make(chan string)
	go s.keepAlivePort(messages)
	message, _ := s.dataUpdatePort()
	// Decode Protocol Buffer & Update everloop Struct LED Count
	proto.Unmarshal([]byte(message), &everloop)
	log.Debug(everloop)
	log.Debug(everloop.EverloopLength)
	// everloop.
	// Print Data
	// log.Debug("Humidity: %f \t Temperature: %f\n", humidity.Humidity, humidity.Temperature)
	messages <- "Data Update Port: CONNECTED"
	// Start Base Port
	// go basePort() // Send Configuration Message

	// Close Data Update Port
	return int(everloop.GetEverloopLength()), everloop.GetLed()

}

func (s *EverloopSensor) initLeds(l int, leds []*core.LedValue) (int, []*core.LedValue) {

	if len(leds) == 0 {
		log.Debug("size 0")
		config := core.EverloopImage{}
		for i := 0; i < l; i++ {
			led := core.LedValue{
				Red:   0,
				Green: 0,
				Blue:  0,
				White: 0,
			}
			// Add New LED to Everloop LED Array
			config.Led = append(config.Led, &led)
		}
		// config.
		return s.SetData(config)
	}
	return l, leds
}
