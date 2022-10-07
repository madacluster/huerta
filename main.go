package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Global Vars
var ZERO_HOST = "tcp://192.168.0.110"
var MOSQUITTO_HOST = "mqtt://192.168.0.108:1883"

func main() {
	fmt.Println("Starting MATRIX CORE Everloop")

	// Asynchronously Start MATRIX CORE Ports
	// go keepAlivePort()  // Ping Once
	// go errorPort()      // Report Any Errors
	// go dataUpdatePort() // Receive # Of Leds & Starts basePort()

	// // Wait For Each Port Connection (ensures each goroutine is able to run)
	// for portStatus := range portStatus {
	// 	fmt.Println("received", portStatus)
	// }
	ZERO_HOST = os.Getenv("ZERO_HOST")
	MOSQUITTO_HOST = os.Getenv("MOSQUITTO_HOST")

	forward()
}
func getHost(port string) string {
	return fmt.Sprintf("%s:%s", ZERO_HOST, port)
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func forward() {
	humidity := NewHumiditySensor()
	uv := NewUvSensor()
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(MOSQUITTO_HOST).SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go humidity.forward(c)
	uv.forward(c)
}
