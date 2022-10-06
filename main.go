package main

import (
	"fmt"

	core "github.com/matrix-io/matrix-protos-go/matrix_io/malos/v1"
)

// Global Vars
var portStatus = make(chan string, 4) // Channel To Ensure Port Goroutines Are Called
var everloop = core.EverloopImage{}   // State Of All MATRIX LEDs
var host = "tcp://192.168.0.110"

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
	webserver()
}
func getHost(port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
