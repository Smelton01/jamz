package ui

import (
	"context"
	"fmt"

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
	devices     list.Model
	playlist    list.Model
	controls    list.Model
	curWindow   window
	controller  control.Controller
	ctx         context.Context
	playlistNav []list.Model
}

func Render(client *spotify.Client) error {

	model := model{controller: control.Controller{Client: client}, ctx: context.Background()}
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
			// m.controls.SetSize(m.devices.Width(), m.devices.Height())
			// m.playlist.SetSize(m.devices.Width(), m.devices.Height())
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
	}

	return mod, cmd
}

func updateDevice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if len(m.devices.Items()) == 0 {
		generic := []interface{}{}
		devs, err := m.controller.GetDevices(context.Background())
		if err != nil {
			panic(err)
		}
		for _, d := range devs {
			generic = append(generic, d)
		}
		m.devices = makeList(generic...)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeWindow(&m.devices, msg)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, nil
		} else if msg.String() == "enter" {
			devs, err := m.controller.GetDevices(context.Background())
			if err != nil {
				panic(err)
			}
			// state, err := m.controller.GetState(m.ctx)
			// if err != nil {
			// 	panic(err)
			// }
			// play on selected device
			for _, dev := range devs {
				if m.devices.SelectedItem().FilterValue() == dev.Name {
					err := m.controller.PlayOpt(m.ctx, &spotify.PlayOptions{DeviceID: &dev.ID})
					if err != nil {
						fmt.Println(err)
						return m, nil
					}
				}
			}
			return m, nil
		}
		var cmd tea.Cmd
		m.devices, cmd = m.devices.Update(msg)

		return m, cmd
	}
	return m, nil
}

func updatePlaylist(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if len(m.controls.Items()) == 0 {
		playlist, err := m.controller.GetPlaylists(context.Background())
		if err != nil {
			panic(err)
		}
		genPlaylist := []interface{}{}
		for _, p := range playlist {
			genPlaylist = append(genPlaylist, p)
		}
		m.playlistNav = append(m.playlistNav, makeList(genPlaylist...))
		m.playlist = m.playlistNav[len(m.playlistNav)-1]
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeWindow(&m.playlist, msg)
	case tea.KeyMsg:
		if len(m.playlistNav) == 1 {
			if msg.String() == "enter" {
				playlist, err := m.controller.GetPlaylists(context.Background())
				if err != nil {
					panic(err)
				}
				for _, p := range playlist {
					if p.Name == m.playlist.SelectedItem().FilterValue() {
						// gen := []interface{}{}
						// for _, track := range p.Tracks {
						// 	gen = append(gen, track)
						// }
					}
				}
				return m, nil
			}

		}
		var cmd tea.Cmd
		m.playlist, cmd = m.playlist.Update(msg)

		return m, cmd
	}
	return m, nil
}

func updateControl(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	if len(m.controls.Items()) == 0 {
		controls := []interface{}{
			item{title: string(playCommand), desc: "resume playback"},
			item{title: string(pauseCommand), desc: "pause playback"},
			item{title: string(nextCommand), desc: "next track"},
			item{title: string(prevCommand), desc: "previous track"},
		}
		m.controls = makeList(controls...)
	}
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
	c := m.controller
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
