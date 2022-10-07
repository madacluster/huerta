package main

import (
	"fmt"
	"net/http"

	pkg "github.com/madacluster/huerto/pkg"
)

const host = "localhost"

// Hudimity
// GET ${host}/${path}/humidity
func Hudimity(w http.ResponseWriter, r *http.Request) {
	sensor := pkg.NewHumiditySensor(host)

	// time.Sleep(1000 * time.Millisecond)
	humidity, temp := sensor.GetData()
	w.Write([]byte(fmt.Sprintf("Humidity: %f\nTemperature: %f", humidity, temp)))
}

// POST ${host}/${path}/temperature/configure

// Uv
// GET ${host}/${path}/uv
func Uv(w http.ResponseWriter, r *http.Request) {
	sensor := pkg.NewUvSensor(host)
	uv, risk := sensor.GetData()
	w.Write([]byte(fmt.Sprintf("Uv: %f\nRisk: %s", uv, risk)))
}

// POST ${host}/${path}/uv/configure

// Pressure
// GET ${host}/${path}/pressure
func Pressure(w http.ResponseWriter, r *http.Request) {
	sensor := pkg.NewPressureSensor(host)
	// time.Sleep(1000 * time.Millisecond)
	pressure, altitude, temp := sensor.GetData()
	w.Write([]byte(fmt.Sprintf("Pressure: %f\nAltitude: %f\nTemperature: %f", pressure, altitude, temp)))
}

// POST ${host}/${path}/pressure/configure

// GPIO
// GET ${host}/${path}/gpio
func Gpio(w http.ResponseWriter, r *http.Request) {
	sensor := pkg.NewGPIOSensor(host)
	// time.Sleep(1000 * time.Millisecond)
	on, off := sensor.GetData()

	w.Write([]byte(fmt.Sprintf("On: %v\nOff: %v\n", on, off)))
}

// POST ${host}/${path}/gpio/${pin}/configure

func webserver() {
	http.HandleFunc("/humidity", Hudimity)
	http.HandleFunc("/pressure", Pressure)
	http.HandleFunc("/uv", Uv)
	http.ListenAndServe(":8080", nil)
}

func main() {
	webserver()
}
