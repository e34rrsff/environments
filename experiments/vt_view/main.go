// These experiments are sloppy work

package main

import (
	_ "charm.land/bubbles/v2"
	"charm.land/bubbletea/v2"
	_ "charm.land/lipgloss/v2"
	"github.com/taigrr/bubbleterm"
	"log"
	_ "log"
	"os/exec"
	_ "strings"
)

type model struct {
	width, height int
	terminal      *bubbleterm.Model
}

func (m *model) Init() tea.Cmd {
	return tea.RequestWindowSize
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.terminal == nil {
			cmd := exec.Command("bash")
			cmd.Env = []string{
				"PATH=/run/current-system/profile/bin",
			}

			t, err := bubbleterm.NewWithCommand(msg.Width, msg.Height, cmd)
			if err != nil {
				log.Fatal(err)
			}
			m.terminal = t
			return m, m.terminal.Init()
		}
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			m.terminal.Close()
			return m, tea.Quit
		}
	}

	if m.terminal != nil {
		terminalModel, cmd := m.terminal.Update(msg)
		m.terminal = terminalModel.(*bubbleterm.Model)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() tea.View {
	if m.terminal == nil {
		return tea.View{Content: "Starting terminal...", AltScreen: true}
	}
	return tea.View{Content: m.terminal.View().Content, AltScreen: true}
}

func main() {
	p := tea.NewProgram(&model{})
	if _, err := p.Run(); err != nil {
		tea.LogToFile("debug.log", "debug")
		log.Fatal(err)
	}
}
