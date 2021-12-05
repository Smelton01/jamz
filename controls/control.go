package control

import (
	"context"

	// "github.com/zmb3/spotify"
	"github.com/zmb3/spotify/v2"
)

type Controller struct {
	Client *spotify.Client
}

// toggle playback
func (c Controller) Play(ctx context.Context) error {
	if err := c.Client.Play(ctx); err != nil {
		return err
	}
	return nil
}
func (c *Controller) Next(ctx context.Context) error {
	if err := c.Client.Next(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Controller) Current(ctx context.Context) (*spotify.FullTrack, error) {
	curr, err := c.Client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return nil, err
	}
	return curr.Item, nil
}

func (c *Controller) Pause(ctx context.Context) error {
	if err := c.Client.Pause(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Controller) Prev(ctx context.Context) error {
	if err := c.Client.Previous(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Controller) Search(ctx context.Context, query string, t spotify.SearchType) (*spotify.SearchResult, error) {
	res, err := c.Client.Search(ctx, query, t)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Controller) GetPlaylists(ctx context.Context) ([]spotify.SimplePlaylist, error) {
	res, err := c.Client.CurrentUsersPlaylists(ctx, spotify.Limit(10))
	if err != nil {
		return nil, err
	}
	// spotify.RequestOption{}
	return res.Playlists, nil
}

func (c *Controller) GetPlaylistTracks(ctx context.Context, ID spotify.ID) ([]spotify.PlaylistTrack, error) {
	res, err := c.Client.GetPlaylistTracks(ctx, ID)
	if err != nil {
		return nil, err
	}
	return res.Tracks, nil
}

func (c *Controller) GetDevices(ctx context.Context) ([]spotify.PlayerDevice, error) {
	devs, err := c.Client.PlayerDevices(ctx)
	if err != nil {
		return nil, err
	}
	return devs, nil
}

func (c *Controller) PlayOpt(ctx context.Context, opt *spotify.PlayOptions) error {
	err := c.Client.PlayOpt(ctx, opt)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetState(ctx context.Context) (*spotify.PlayerState, error) {
	res, err := c.Client.PlayerState(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// func (c *Controller) GetCurrentTrack
