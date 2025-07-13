/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/barnardn/go-clime/geolocatedio"
	"github.com/barnardn/go-clime/ipclient"
	"github.com/barnardn/go-clime/whirly"
	"github.com/spf13/cobra"
)

// geocodeCmd represents the geocode command
var geocodeCmd = &cobra.Command{
	Use:   "geocode",
	Short: "Fetch geographic details based on ip address",
	Run:   runGeocode,
}

var (
	justZip bool
)

func init() {
	rootCmd.AddCommand(geocodeCmd)
	geocodeCmd.Flags().BoolVarP(&justZip, "zipcode", "z", false, "Display only the zip code")
}

func runGeocode(cmd *cobra.Command, args []string) {
	progress := whirly.New(whirly.Kitt)
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
	fmt.Printf("Zip Code: %s\n", location.ZipCode)
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
	geoClient := geolocatedio.NewClient("<your api key here>")
	geoChan, errChan := geoClient.GeoLocation(ipAddress)
	select {
	case err := <-errChan:
		return nil, err
	case geo := <-geoChan:
		return &geo, nil
	}

}
