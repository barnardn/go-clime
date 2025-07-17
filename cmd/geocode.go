/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/barnardn/go-clime/geolocatedio"
	"github.com/barnardn/go-clime/ipclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// geocodeCmd represents the geocode command
var geocodeCmd = &cobra.Command{
	Use:   "geocode",
	Short: "Fetch geographic details based on ip address",
	Run:   runGeocode,
}

func init() {
	rootCmd.AddCommand(geocodeCmd)
	geocodeCmd.Flags().BoolVarP(&justZip, "zipcode", "z", false, "Display only the zip code")
	geocodeCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "Quiet mode. Don't show progress indicator")
}

func runGeocode(cmd *cobra.Command, args []string) {
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
	progress.Stop()
	if justZip {
		fmt.Println(location.ZipCode)
	} else {
		fmt.Println(location.String())
	}
}

func runIPfetch() (*string, error) {
	ipClient := ipclient.NewClient()
	ipChan, errChan := ipClient.PublicIP()
	select {
	case err := <-errChan:
		return nil, err
	case ipAddress := <-ipChan:
		return &ipAddress, nil
	}
}

func geocode(ipAddress string) (*geolocatedio.LocationInfo, error) {
	geoClient := geolocatedio.NewClient(viper.GetString("geokey"))
	geoChan, errChan := geoClient.GeoLocation(ipAddress)
	select {
	case err := <-errChan:
		return nil, err
	case geo := <-geoChan:
		return &geo, nil
	}
}
