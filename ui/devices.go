package ui

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
	dev  chan list.Item
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.dev <- list.Model{}
			return m, nil
		} else if msg.String() == "enter" {
			// log.Println(m.list.SelectedItem())
			m.dev <- m.list.SelectedItem()
			log.Println("sent something")
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Main(devs []spotify.PlayerDevice, ch chan list.Item) {
	items := []list.Item{}

	for _, dev := range devs {
		items = append(items, item{
			title: dev.Name, desc: dev.Type,
		})
	}
	if len(items) < 1 {
		items = append(items, item{title: "No device detected", desc: "Please make sure your device is connected to the internet"})
	}

	ch1 := make(chan list.Item)
	m := model{list: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), dev: ch1}
	m.list.Title = "JAMZ select Active device"

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	go func(ch, ch1 chan list.Item) {
		dev := <-ch1
		log.Println("got", dev)
		ch <- dev

	}(ch, ch1)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		// os.Exit(1)
	}
	log.Println()
}
