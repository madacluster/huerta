package main

import (
	"fmt"

	pkg "github.com/madacluster/huerto/pkg"
	"github.com/spf13/cobra"
)

var ZERO_HOST = "tcp://192.168.0.110"

func main() {
	var echoTimes int
	var cmdHumidity = &cobra.Command{
		Use:   "humidity",
		Short: "Print humidity",
		Long:  `print humidity and temperature from sensor`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewHumiditySensor(ZERO_HOST)
			humidity, temp := sensor.GetData()
			fmt.Printf("Humidity: %f\nTemperature: %f\n", humidity, temp)
		},
	}

	var cmdPressure = &cobra.Command{
		Use:   "pressure",
		Short: "Print Pressure",
		Long:  `Print pressure, altitude and temperature from sensor`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewPressureSensor(ZERO_HOST)
			pressure, altitude, temp := sensor.GetData()
			fmt.Printf("Pressure: %f\nAltitude: %f\nTemperature: %f", pressure, altitude, temp)
		},
	}

	var cmdUV = &cobra.Command{
		Use:   "uv",
		Short: "Print UV",
		Long:  `Print UV and risk from sensor`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewUvSensor(ZERO_HOST)
			index, risk := sensor.GetData()
			fmt.Printf("index: %f\nrisk: %s\n", index, risk)
		},
	}

	var cmdGPIO = &cobra.Command{
		Use:   "gpio",
		Short: "Print GPIO",
		Long:  `Print GPIO pin and value from sensor`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewGPIOSensor(ZERO_HOST)
			pin, mode, value, values := sensor.GetData()
			fmt.Printf("Pin: %d\nMode: %v\nValue: %d\nValues: %d\n", pin, mode, value, values)
		},
	}

	cmdUV.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "matrix"}
	rootCmd.AddCommand(cmdHumidity, cmdPressure, cmdUV, cmdGPIO)
	cmdPressure.AddCommand(cmdUV)
	rootCmd.Execute()
}
