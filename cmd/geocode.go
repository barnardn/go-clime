/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/barnardn/go-clime/ipclient"
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
	ipClient := ipclient.NewClient()
	ipAddress, err := ipClient.FetchIP()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Print(*ipAddress)
}
