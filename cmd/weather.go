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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Returns the current weather conditions at the provide zip code",
	Run:   runWeather,
}

var weatherHereCmd = &cobra.Command{
	Use:   "here",
	Short: "Returns the weather conditions at the current location",
	Run:   runWeatherHere,
}

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.Flags().BoolVarP(&isImperial, "imperial", "i", false, "Display imperial measurements")
	weatherCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "Quiet mode. Don't show progress indicator")
	weatherCmd.Flags().BoolVarP(&justTemp, "just-temp", "t", false, "Return only the temperature")
	weatherCmd.Flags().BoolVarP(&justFeelsLike, "just-feelslike", "f", false, "Return only the feels like temperature")

	rootCmd.AddCommand(weatherHereCmd)
	weatherHereCmd.Flags().BoolVarP(&isImperial, "imperial", "i", false, "Display imperial measurements")
	weatherHereCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "Quiet mode. Don't show progress indicator")
	weatherHereCmd.Flags().BoolVarP(&justTemp, "just-temp", "t", false, "Return only the temperature")
	weatherHereCmd.Flags().BoolVarP(&justFeelsLike, "just-feelslike", "f", false, "Return only the feels like temperature")
}

func runWeatherHere(cmd *cobra.Command, args []string) {
	progress := newWhirly()
	progress.Start()

	ipAddress, err := runIPfetch()
	if err != nil {
		progress.Stop()
		log.Fatalf("%+v\n", err)
	}
	location, err := geocode(*ipAddress)
	if err != nil {
		progress.Stop()
		log.Fatalf("%+v\n", err)
	}

	units := clime.Metric
	if isImperial {
		units = clime.Imperial
	}
	climeClient := clime.NewClient(
		openweathermap.NewClient(viper.GetString("apikey"), isImperial),
	)

	ccChan, errChan := climeClient.AsyncConditions(location.ZipCode)
	select {
	case err := <-errChan:
		progress.Stop()
		log.Fatalf("%+v\n", err)
	case conditions := <-ccChan:
		progress.Stop()
		cc := clime.NewCurrentConditions(conditions, units)
		if justTemp {
			fmt.Println(cc.TemperatureDetails.Current.String())
		} else if justFeelsLike {
			fmt.Println(cc.TemperatureDetails.FeelsLike.String())
		} else {
			fmt.Print(cc.String())
		}
	}

}

func runWeather(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}
	progress := newWhirly()
	progress.Start()

	cc, err := currentConditions(args[0])
	progress.Stop()

	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if justTemp {
		fmt.Println(cc.TemperatureDetails.Current.String())
	} else if justFeelsLike {
		fmt.Println(cc.TemperatureDetails.FeelsLike.String())
	} else {
		fmt.Print(cc.String())
	}
}

func currentConditions(zipCode string) (*clime.CurrentConditions, error) {
	units := clime.Metric
	if isImperial {
		units = clime.Imperial
	}
	climeClient := clime.NewClient(
		openweathermap.NewClient(viper.GetString("apikey"), isImperial),
	)
	ccChan, errChan := climeClient.AsyncConditions(zipCode)
	select {
	case err := <-errChan:
		return nil, err
	case conditions := <-ccChan:
		cc := clime.NewCurrentConditions(conditions, units)
		return &cc, nil
	}
}
