package ui

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	control "github.com/smelton01/jamz/controls"
	"github.com/zmb3/spotify/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type window int
type command string

const (
	deviceWindow window = iota
	playlistWindow
	controlWindow
	nowPlayingWindow
	searchWindow
	recommendedWindow
)

const (
	playCommand  command = "play"
	pauseCommand command = "pause"
	nextCommand  command = "next"
	prevCommand  command = "prev"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	devices   list.Model
	playlist  list.Model
	controls  list.Model
	curWindow window
	client    *spotify.Client
}

func Main(devs []spotify.PlayerDevice, client *spotify.Client) {
	items := []list.Item{}

	for _, dev := range devs {
		items = append(items, item{
			title: dev.Name, desc: dev.Type,
		})
	}
	if len(items) < 1 {
		items = append(items, item{title: "No device detected", desc: "Please make sure your device is connected to the internet"})
	}
	playlist := []list.Item{
		item{title: "first", desc: "etsting"},
		item{title: "2nd", desc: "second one"}}
	controls := []list.Item{
		item{title: string(playCommand), desc: "resume playback"},
		item{title: string(pauseCommand), desc: "pause playback"},
		item{title: string(nextCommand), desc: "next track"},
		item{title: string(prevCommand), desc: "previous track"},
	}
	m := model{devices: list.NewModel(items, list.NewDefaultDelegate(), 0, 0),
		playlist: list.NewModel(playlist, list.NewDefaultDelegate(), 0, 0),
		controls: list.NewModel(controls, list.NewDefaultDelegate(), 0, 0),
		client:   client,
	}
	m.devices.Title = "JAMZ select Active device"
	m.playlist.Title = "Playlists"
	m.controls.Title = "Playback controls"

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	log.Println()
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msgStr := msg.String(); msgStr == "ctrl+c" || msgStr == "q" {
			return nil, tea.Quit
		} else if msg.String() == "1" {
			m.curWindow = (m.curWindow + 1) % 3
		}
	}
	var cmd tea.Cmd
	var mod tea.Model
	switch m.curWindow {
	case deviceWindow:
		mod, cmd = updateDevice(msg, m)
	case playlistWindow:
		mod, cmd = updatePlaylist(msg, m)
	case controlWindow:
		mod, cmd = updateControl(msg, m)

		// return m, cmd
	}

	return mod, cmd
}

func updateDevice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.devices.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			// m.dev <- list.Model{}
			return m, nil
		} else if msg.String() == "enter" {
			// m.dev <- m.list.SelectedItem()
			log.Println("selected", m.devices.SelectedItem())
			return m, nil
		}
		var cmd tea.Cmd
		m.devices, cmd = m.devices.Update(msg)

		return m, cmd
	}
	return m, nil
}

func updatePlaylist(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.playlist.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// m.dev <- m.list.SelectedItem()
			log.Println("selected playlist")
			return m, nil
		}
		var cmd tea.Cmd
		m.playlist, cmd = m.playlist.Update(msg)

		return m, cmd
	}
	return m, nil
}

func updateControl(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.controls.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// m.dev <- m.list.SelectedItem()
			err := ctl(context.Background(), m.controls.SelectedItem().FilterValue(), m)
			if err != nil {
				// fmt.Println(err)
				// Show a thing if this fails maybe red color or something
				return m, nil
			}
			// fmt.Println("enter")
			return m, nil
		}
		var cmd tea.Cmd
		m.controls, cmd = m.controls.Update(msg)
		return m, cmd
	}
	return m, nil

}

func ctl(ctx context.Context, cmd string, m model) error {
	c := control.Controller{Client: m.client}
	// fmt.Println(c)
	// err := c.Client.Play(ctx)
	// if err != nil {
	// 	fmt.Println("got errsssssssssssssssssssssssssssssssssssssssor", err)
	// 	return err
	// }
	switch command(cmd) {
	case playCommand:
		return c.Play(ctx)
	case pauseCommand:
		return c.Pause(ctx)
	case nextCommand:
		return c.Next(ctx)
	case prevCommand:
		return c.Prev(ctx)

	}
	return nil
}

func (m model) View() string {
	var view string
	switch m.curWindow {
	case deviceWindow:
		view = m.devices.View()
	case playlistWindow:
		view = m.playlist.View()
	case controlWindow:
		view = m.controls.View()
	}
	return docStyle.Render(view)
}
