package main

import (
	"fmt"
	"strconv"

	pkg "github.com/madacluster/huerto/pkg"
	"github.com/spf13/cobra"
)

func Humidity() *cobra.Command {
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

	calibrate := &cobra.Command{
		Use:   "calibrate [temp]",
		Short: "Calibrate temperatue",
		Long:  `calibrate temperature from sensor`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No temperature specified")
				return
			}
			temp, _ := strconv.ParseFloat(args[0], 32)
			sensor := pkg.NewHumiditySensor(ZERO_HOST)
			humidity, ctemp := sensor.Calibrate(float32(temp))
			fmt.Printf("Humidity: %f\nTemperature: %f\n", humidity, ctemp)
		},
	}
	cmdHumidity.AddCommand(calibrate)
	return cmdHumidity

}
