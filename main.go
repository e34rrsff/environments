package main

import (
	"log"
	"os"
	"slices"
	"strings"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alexflint/go-arg"
)

var version = "dev"

var style = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Blue)

func init() {
	var args args
	arg.MustParse(&args)
}

func main() {
	m := mainModel{}
	m.addView(&menuView{items: []string{
		"apple",
		"banana",
		"orange",
	}})

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

// Argument parsing stuff

type args struct {
	DummyArg string `arg:",positional"`
}

func (args) Version() string {
	return version
}

// Messages

type deleteViewMsg struct {
	viewToDelete view
}

type addViewMsg struct {
	viewToAdd view
	// have to pass width and height as a new view is being created since
	// the main model does not keep track of this
	w, h int
}

// Boilerplate handling for the "views"

//	"Views" (not the function) will be both a compatible `view` struct and
//	`tea.Model`. They require width and height fields, but other than that,
//	they can contain whatever fields want!

type sizable struct {
	width, height int
}

func (s *sizable) setSize(width, height int) {
	s.width, s.height = width, height
}

type view interface {
	tea.Model
	setSize(width, height int)
}

// for lack of a better name, `liveView`s will be views that should have a
// dedicated function to exit cleany  (i.g. a terminal view)
type liveView struct {
	cleanUp func()
}

// the main "underlying" model serving as a preliminary event handler
type mainModel struct {
	stack []view
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) View() tea.View {
	if len(m.stack) == 0 {
		return tea.View{Content: ""}
	}

	top := m.stack[len(m.stack)-1]
	styled := style.Render(top.View().Content)

	return tea.View{
		Content:   styled,
		AltScreen: true,
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.stack) == 0 {
		return m, tea.Quit
	}

	top := m.stack[len(m.stack)-1]
	_, cmd := top.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top.setSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+alt+q":
			return m, tea.Quit // TODO handle quit more elegantly
		}

	case addViewMsg:
		msg.viewToAdd.setSize(msg.w, msg.h)
		m.stack = append(m.stack, msg.viewToAdd)

	case deleteViewMsg:
		for i, view := range m.stack {
			if msg.viewToDelete == view {
				m.stack = slices.Delete(m.stack, i, i+1)
			}
		}
	}

	return m, cmd
}

func (m *mainModel) addView(v view) {
	m.stack = append(m.stack, v)
}

// selectable items list view
type menuView struct {
	sizable
	items  []string
	cursor int
}

func (m menuView) Init() tea.Cmd {
	return nil
}

func (m menuView) View() tea.View {
	var content strings.Builder
	content.WriteString("available exams:\n\n")
	for _, item := range m.items {
		if m.items[m.cursor] == item {
			content.WriteString("\u2192")
		}
		content.WriteString(" ")
		content.WriteString(item)
		content.WriteString("\n")
	}

	formatted := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content.String())

	return tea.View{Content: formatted}
}

func (m *menuView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor != 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor != len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			selection := m.items[m.cursor]
			return m, func() tea.Msg {
				return addViewMsg{
					viewToAdd: &selectionView{message: selection},
					w:         m.width,
					h:         m.height,
				}
			}
		}
	}

	return m, nil
}

// simple view to display a centered error message
type selectionView struct {
	sizable
	message string
}

func (m selectionView) Init() tea.Cmd {
	return nil
}

func (m *selectionView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyPressMsg:
		return m, func() tea.Msg { return deleteViewMsg{viewToDelete: m} }
	}
	return m, nil
}

func (m *selectionView) View() tea.View {
	content := "You selected:\n\n" + m.message
	return tea.View{
		Content: lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)}
}
