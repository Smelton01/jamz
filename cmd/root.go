/*
Copyright © 2021 Simon Juba

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/smelton01/jamz/login"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"

	"github.com/spf13/viper"
)

var (
	cfgFile string
)

const (
	id     = "SPOTIFY_ID"
	secret = "SPOTIFY_SECRET"
)

var Client *spotify.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spotui",
	Short: "A Terminal based Spotify app",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
d`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: nil,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Spotify......")
		acc := login.MakeAcc()
		client, err := acc.Auth()
		if err != nil {
			panic(err)
		}
		Client = client

		// saveToken(client)
		// tracks, err := client.PlayerRecentlyPlayed(acc.Ctx)
		// if err != nil {
		// 	panic(err)
		// }
		// for _, track := range tracks {
		// 	fmt.Printf("Name: %v, Artist: %v", track.Track.Name, track.Track.Artists)
		// }

		dev, err := Client.PlayerDevices(context.Background())
		if err != nil {
			panic(err)
		}
		for _, d := range dev {
			fmt.Printf("d.Name: %v  ", d.Name)
			fmt.Printf("d.ID: %v\n", d.ID.String())
			client.PlayOpt(cmd.Context(), &spotify.PlayOptions{
				DeviceID: &d.ID,
			})
			client.Play(cmd.Context())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spotui.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".spotui" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".credentials")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	os.Setenv(id, viper.GetString(id))
	os.Setenv(secret, viper.GetString(secret))

	acc := login.MakeAcc()
	client, err := acc.Auth()
	if err != nil {
		panic(err)
	}
	Client = client
}
