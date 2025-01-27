/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/barnardn/go-clime/clime"
	"github.com/barnardn/go-clime/openweathermap"
	"github.com/barnardn/go-clime/whirly"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Returns the current weather conditions",
	Run:   runWeather,
}

var (
	isImperial bool
	isQuiet    bool
)

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.Flags().BoolVarP(&isImperial, "imperial", "i", false, "Display imperial measurements")
	weatherCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "Quiet mode. Don't show progress indicator")
}

func runWeather(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}
	units := clime.Metric
	if isImperial {
		units = clime.Imperial
	}
	climeClient := clime.NewClient(
		openweathermap.NewClient(viper.GetString("apikey"), isImperial),
	)
	whirlyType := whirly.Kitt
	if isQuiet {
		whirlyType = whirly.Empty
	}
	progress := whirly.New(whirlyType)
	progress.Start()

	ccChan, errChan := climeClient.AsyncConditions(args[0])
	select {
	case err := <-errChan:
		progress.Stop()
		log.Fatalf("%+v\n", err)
	case conditions := <-ccChan:
		progress.Stop()
		cc := clime.NewCurrentConditions(conditions, units)
		fmt.Print(cc.String())
	}
}
