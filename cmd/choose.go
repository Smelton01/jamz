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
	Use:   "old",
	Short: "A brief description of your command",
	Long:  `Select device to play on`,
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
		for _, dev := range devs {
			fmt.Println(dev.Name, "status", dev.Active)
			if dev.Name == device.FilterValue() {
				dev.Active = true
				playID = dev.ID
			}
			fmt.Println(dev.Name, "status", dev.Active)
		}
		err = Client.PlayOpt(context.Background(), &spotify.PlayOptions{DeviceID: &playID})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(oldCmd)
}
