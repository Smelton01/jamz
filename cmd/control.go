package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play music",
	Long:  `Resume playback on yoour currenlyt active device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return play(context.Background())
	},
}

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause music",
	Long:  `Pause playback on your currently active device`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := pause(context.Background())
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
		err := next(cmd.Context())
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
		err := prev(cmd.Context())
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
		results, err := search(cmd.Context(), query, sType)
		if err != nil {
			return err
		}

		/// TODO use the results to display somehting
		_ = results
		return nil
	},
}

// toggle playback
func play(ctx context.Context) error {
	// state, err := Client.PlayerState(ctx)
	// if err != nil {
	// 	return err
	// }

	// if state.Device == spotify.PlayerDevice{}{
	// 	return ErrNoActiveDevice
	// }
	// // if dev := state.Device; dev == spotify.PlayerDevice{} {
	// // 	return ErrNoActiveDevice
	// // }

	if err := Client.Play(ctx); err != nil {
		return err
	}
	return nil
}
func next(ctx context.Context) error {
	if err := Client.Next(ctx); err != nil {
		return err
	}
	return nil
}

func pause(ctx context.Context) error {
	if err := Client.Pause(ctx); err != nil {
		return err
	}
	return nil
}

func prev(ctx context.Context) error {
	if err := Client.Previous(ctx); err != nil {
		return err
	}
	return nil
}

func search(ctx context.Context, query string, t spotify.SearchType) (*spotify.SearchResult, error) {
	res, err := Client.Search(ctx, query, t)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func init() {
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(prevCmd)
	rootCmd.AddCommand(pauseCmd)

}
