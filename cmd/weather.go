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

func init() {
	rootCmd.AddCommand(weatherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weatherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weatherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runWeather(cmd *cobra.Command, args []string) {
	fmt.Println("run weather")
	clime := clime.NewClient(
		openweathermap.NewClient(viper.GetString("apikey")),
	)
	conditions, err := clime.CurrentConditions("49002")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Printf("%+v\n", conditions)
}
