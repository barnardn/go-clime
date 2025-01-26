/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/barnardn/go-clime/clime"
	"github.com/barnardn/go-clime/openweathermap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Returns the current weather conditions",
	Run:   runWeather,
}

var isImperial bool

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.Flags().BoolVarP(&isImperial, "imperial", "i", false, "Display imperial measurements")
}

func runWeather(cmd *cobra.Command, args []string) {
	units := clime.Metric
	if isImperial {
		units = clime.Imperial
	}
	climeClient := clime.NewClient(
		openweathermap.NewClient(viper.GetString("apikey"), isImperial),
	)
	conditions, err := climeClient.CurrentConditions(args[0])
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	cc := clime.NewCurrentConditions(*conditions, units)
	fmt.Print(cc.String())
}
