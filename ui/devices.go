package ui

import (
	"context"
	"fmt"
	"log"

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
	devices    list.Model
	playlist   list.Model
	controls   list.Model
	curWindow  window
	client     *spotify.Client
	controller control.Controller
}

func Render(client *spotify.Client) error {

	model := model{controller: control.Controller{Client: client}}

	generic := []interface{}{}
	devs, err := model.controller.GetDevices(context.Background())
	if err != nil {
		return err
	}
	for _, d := range devs {
		generic = append(generic, d)
	}
	devices := makeList(generic...)
	playlist := []list.Item{
		item{title: "first", desc: "etsting"},
		item{title: "2nd", desc: "second one"}}
	genPlaylist := []interface{}{}
	for _, p := range playlist {
		genPlaylist = append(genPlaylist, p)
	}

	controls := []interface{}{
		item{title: string(playCommand), desc: "resume playback"},
		item{title: string(pauseCommand), desc: "pause playback"},
		item{title: string(nextCommand), desc: "next track"},
		item{title: string(prevCommand), desc: "previous track"},
	}

	model.devices = devices
	model.playlist = makeList(genPlaylist...)
	model.controls = makeList(controls...)
	model.devices.Title = "JAMZ select Active device"
	model.playlist.Title = "Playlists"
	model.controls.Title = "Playback controls"

	p := tea.NewProgram(model)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		return err
	}
	return nil
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
		resizeWindow(&m.devices, msg)
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
		resizeWindow(&m.playlist, msg)
	case tea.KeyMsg:
		if msg.String() == "enter" {
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
		resizeWindow(&m.controls, msg)
	case tea.KeyMsg:
		if msg.String() == "enter" {
			err := ctl(context.Background(), m.controls.SelectedItem().FilterValue(), m)
			if err != nil {
				// fmt.Println(err)
				// Show a thing if this fails maybe red color or something
				return m, nil
			}
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
func makeList(items ...interface{}) list.Model {
	output := []list.Item{}
	for _, elem := range items {
		switch elemType := elem.(type) {
		case spotify.PlayerDevice:
			output = append(output, item{title: elemType.Name, desc: elemType.Type})
			if len(output) < 1 {
				output = append(output, item{title: "No device detected", desc: "Please make sure your device is connected to the internet"})
			}
		case item:
			output = append(output, item{title: elemType.title, desc: elemType.desc})
		case spotify.SimplePlaylist:
			output = append(output, item{title: elemType.Name, desc: elemType.Owner.DisplayName})
		}

	}
	return list.NewModel(output, list.NewDefaultDelegate(), 0, 0)
}

func resizeWindow(list *list.Model, msg tea.WindowSizeMsg) {
	top, right, bottom, left := docStyle.GetMargin()
	list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
}
