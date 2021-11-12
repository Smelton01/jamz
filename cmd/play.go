/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/smelton01/jamz/ui"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// playCmd represents the play command
var oldCmd = &cobra.Command{
	Use:   "play",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		devs, err := Client.PlayerDevices(cmd.Context())
		if err != nil {
			panic(err)
		}
		log.Println(devs)
		ch := make(chan list.Item)
		ui.Main(devs, ch)
		device := <-ch
		// check device
		var playID spotify.ID
		log.Println("Active device")
		for _, dev := range devs {
			fmt.Println(dev.Name, "status", dev.Active)
			if dev.Name == device.FilterValue() {
				dev.Active = true
				playID = dev.ID
			}
			fmt.Println(dev.Name, "status", dev.Active)
		}

		// ui.Main(devs, ch)
		state, err := Client.PlayerState(cmd.Context())
		if err != nil {
			panic(err)
		}
		fmt.Println("device", state.Device, state.Playing, state.Timestamp)
		rec, err := Client.PlayerRecentlyPlayed(context.Background())
		if err != nil {
			panic(err)
		}
		log.Println(rec[0].PlaybackContext.URI)

		err = Client.PlayOpt(context.Background(), &spotify.PlayOptions{DeviceID: &playID})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(oldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
