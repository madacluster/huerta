package main

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	pkg "github.com/madacluster/huerto/pkg"
	"github.com/spf13/cobra"
)

var ZERO_HOST = "tcp://"
var debug bool

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func toggleDebug(cmd *cobra.Command, args []string) {
	if debug {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else {
		plainFormatter := new(PlainFormatter)
		log.SetFormatter(plainFormatter)
	}
}
func main() {
	var cmdHumidity = Humidity()

	var cmdPressure = &cobra.Command{
		Use:   "pressure",
		Short: "Print Pressure",
		Long:  `Print pressure, altitude and temperature from sensor`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewPressureSensor(ZERO_HOST)
			pressure, altitude, temp := sensor.GetData()
			fmt.Printf("Pressure: %f\nAltitude: %f\nTemperature: %f\n", pressure, altitude, temp)
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
		Use:    "gpio",
		Short:  "Print GPIO",
		Long:   `Print GPIO pin and value from sensor`,
		PreRun: toggleDebug,
		Args:   cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			sensor := pkg.NewGPIOSensor(ZERO_HOST)
			on, off := sensor.GetData()
			fmt.Printf("On: %v\nOff: %v\n", on, off)
		},
	}

	var cmdReadGPIO = &cobra.Command{
		Use:    "read",
		Short:  "Read Pin GPIO",
		Long:   `Read Pin GPIO`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: toggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No pin specified")
				return
			} else {
				pinIn, _ := strconv.Atoi(args[0])
				sensor := pkg.NewGPIOSensor(ZERO_HOST)
				if sensor.ReadPin(pinIn) {
					fmt.Printf("Pin %d is on\n", pinIn)
				} else {
					fmt.Printf("Pin %d is off\n", pinIn)
				}
			}
		},
	}

	var cmdWriteGPIO = &cobra.Command{
		Use:    "write",
		Short:  "Write Pin GPIO",
		Long:   `Write Pin GPIO`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: toggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No pin specified")
				return
			} else {
				pinIn, err := strconv.Atoi(args[0])
				if err != nil {
					log.Error(err)
				}
				if pinIn > 15 {
					err := fmt.Errorf("max index(15) size < %d", pinIn)
					log.Error("Error: ", err)
					return
				}

				input, err := strconv.Atoi(args[1])
				if err != nil {
					log.Error(err)
				}
				sensor := pkg.NewGPIOSensor(ZERO_HOST)
				sensor.WritePin(pinIn, input)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "matrix"}
	rootCmd.AddCommand(cmdHumidity, cmdPressure, cmdUV, cmdGPIO, Everloop(), Everloop())
	cmdGPIO.AddCommand(cmdReadGPIO, cmdWriteGPIO)
	rootCmd.PersistentFlags().StringVarP(&ZERO_HOST, "host", "H", "tcp://192.168.0.110", "ZeroMQ host")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "verbose logging")

	rootCmd.Execute()
}
