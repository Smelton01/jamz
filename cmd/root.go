/*
Copyright © 2021 Simon Juba
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/smelton01/jamz/login"
	"github.com/smelton01/jamz/ui"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"

	"github.com/spf13/viper"
)

var (
	cfgFile string
	Client  *spotify.Client
)

const (
	id     = "SPOTIFY_ID"
	secret = "SPOTIFY_SECRET"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jamz",
	Short: "Start the TUI",
	Long: `A Terminal based interface for the Spotify API
	The base command fires up the TUI which has multiple windows to navigate
	and keybindings to help navigate and choose the perfect music`,
	PersistentPreRun: nil,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting Spotify......")
		return ui.Render(Client)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spotui.yaml)")

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
	} else {
		os.Setenv(id, viper.GetString(id))
		os.Setenv(secret, viper.GetString(secret))
	}

	acc := login.MakeAcc()
	client, err := acc.Auth()
	if err != nil {
		panic(err)
	}
	Client = client
}
