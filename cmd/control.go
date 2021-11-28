package cmd

import (
	"context"
	"strings"

	control "github.com/smelton01/jamz/controls"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play music",
	Long:  `Resume playback on yoour currenlyt active device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &control.Controller{Client: Client}
		return c.Play(context.Background())
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause music",
	Long:  `Pause playback on your currently active device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &control.Controller{Client: Client}
		err := c.Pause(context.Background())
		if err != nil {
			return err
		}
		return nil
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Skip to the next track",
	Long:  `Next track`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &control.Controller{Client: Client}
		err := c.Next(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
}
var prevCmd = &cobra.Command{
	Use:   "prev",
	Short: "Skip to the previous track",
	Long:  `Prev track`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &control.Controller{Client: Client}
		err := c.Prev(cmd.Context())
		if err != nil {
			return err
		}
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search for something on spotify",
	Long:  `search command with flags`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		query = strings.TrimSpace(query)
		if query == "" {
			return ErrInvalidQuery
		}
		sType := spotify.SearchTypeAlbum
		c := &control.Controller{Client: Client}
		results, err := c.Search(cmd.Context(), query, sType)
		if err != nil {
			return err
		}

		/// TODO use the results to display somehting
		_ = results
		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(pauseCmd)

}
