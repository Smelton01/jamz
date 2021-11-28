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
