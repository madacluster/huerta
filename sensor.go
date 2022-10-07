package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
	zmq "github.com/pebbe/zmq4"
)

type Ports struct {
	// Connect ZMQ Socket To MATRIX CORE
	config string
	keep   string
	error  string
	update string
}

type Sensor struct {
	ports Ports
	name  string
}

func NewSensor(name, config, keep, error, update string) *Sensor {
	return &Sensor{
		ports: Ports{
			config,
			keep,
			error,
			update,
		},
		name: name,
	}
}

// BASE PORT \\ (port where configurations are sent)
func (s *Sensor) basePort() {
	// Connect ZMQ Socket To MATRIX CORE
	// zmq.PUSH
	pusher, _ := zmq.NewSocket(zmq.PUSH)    // Create A Pusher Socket
	pusher.Connect(getHost(s.ports.config)) // Connect Pusher To Base Port

	// Notify That Port Is Ready

	// Keep Sending Everloop Image
	for {
		// Create (x) Amount Of Randomly Colored LEDs

		// Create Everloop Driver Configuration Protocol
		configuration := core.DriverConfig{
			Humidity: &core.HumidityParams{
				CurrentTemperature: 20.7,
			},
			DelayBetweenUpdates:  1,
			TimeoutAfterLastPing: 6,
		}
		//Encode Protocol Buffer
		var encodedConfiguration, _ = proto.Marshal(&configuration)
		// Send Protocol Buffer
		pusher.Send(string((encodedConfiguration)), 1)

		// Reset Everloop Array
		// Loop delay
		time.Sleep(50 * time.Millisecond)
	}
}

// KEEP-ALIVE PORT \\ (port where pings are sent)
func (s *Sensor) keepAlivePort(channel chan string) {
	// Connect ZMQ Socket To MATRIX CORE
	pusher, _ := zmq.NewSocket(zmq.PUSH)  // Create A Pusher Socket
	pusher.Connect(getHost(s.ports.keep)) // Connect Pusher To Keep-Alive Port

	// Notify That Port Is Ready

	// Send Keep Alive Message
	for {
		select {
		case <-channel:
			return
		default:
			pusher.Send("", 1)
			fmt.Printf("Keep-Alive Sent for Sensor %s\n", s.name)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

// ERROR PORT \\ (port where errors are received)
func (s *Sensor) errorPort(channel chan string) {
	// Connect ZMQ Socket To MATRIX CORE
	subscriber, _ := zmq.NewSocket(zmq.SUB)    // Create A Subscriber Socket
	subscriber.Connect(getHost(s.ports.error)) // Connect Subscriber To Data Update Port
	subscriber.SetSubscribe("")                // Subscribe To Error Port Messages

	// Notify That Port Is Ready

	// Wait For Error
	for {
		// On Error
		message, _ := subscriber.Recv(2)
		// Print Error
		fmt.Println("ERROR:", message)
	}
}

func (s *Sensor) dataUpdatePort() (string, error) {
	// Connect ZMQ Socket To MATRIX CORE

	subscriber, _ := zmq.NewSocket(zmq.SUB)     // Create A Subscriber Socket
	subscriber.Connect(getHost(s.ports.update)) // Connect Subscriber To Data Update Port
	subscriber.SetSubscribe("")                 // Subscribe To Data Update Port Messages

	// Notify That Port Is Ready

	// Wait For Data
	// On Data
	return subscriber.Recv(2)

}
