package main

import (
	"fmt"
	"strconv"
	"strings"

	pkg "github.com/madacluster/huerto/pkg"
	"github.com/spf13/cobra"
)

func Everloop() *cobra.Command {
	var cmd = &cobra.Command{
		Use:    "everloop",
		Short:  "Print everloop",
		Long:   `Print everloop pin and value from sensor`,
		PreRun: toggleDebug,
		Args:   cobra.MinimumNArgs(1),
		// Run: func(cmd *cobra.Command, args []string) {
		// 	if len(args) == 0 {
		// 		fmt.Println("No pin specified")
		// 		return
		// 	}

		// },
	}
	set := &cobra.Command{
		Use:    "set [led] [red,blue,green,white]",
		Short:  "Set LED",
		Long:   `Write Pin GPIO`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: toggleDebug,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No led specified")
				return
			} else {

				leds := SetLeds(args[0], args[1])
				sensor := pkg.NewEverloopSensor(ZERO_HOST)
				sensor.WriteLed(leds)
			}
		},
	}
	cmd.AddCommand(set)

	return cmd
}

func GetLeds(ls string) []int {
	result := make([]int, 0)
	leds := strings.Split(ls, ",")
	for _, l := range leds {
		intL, _ := strconv.Atoi(l)
		result = append(result, intL)
	}
	return result
}

func GetColor(color string) []int {
	result := make([]int, 4)
	leds := strings.Split(color, ",")
	for i := range result {
		intL, _ := strconv.Atoi(leds[i])
		result[i] = intL
	}
	return result
}

func SetLeds(leds, color string) []*pkg.LED {
	l := GetLeds(leds)
	c := GetColor(color)
	ledList := make([]*pkg.LED, 0)
	for i := 0; i < 35; i++ {
		led := pkg.LED{
			Red:   0,
			Green: 0,
			Blue:  0,
			White: 0,
		}
		// red,blue,green,white
		if inArray(i, l) {
			led.Red = c[0]
			led.Blue = c[1]
			led.Blue = c[2]
			led.White = c[3]
		}
		ledList = append(ledList, &led)
	}
	return ledList
}

func inArray(id int, leds []int) bool {
	for _, l := range leds {
		if l == id {
			return true
		}
	}
	return false
}
