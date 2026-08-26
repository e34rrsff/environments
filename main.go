package main

import (
	_ "embed"
	"log"
	"os"
	"strings"

	"go.yaml.in/yaml/v4"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Argument parsing stuff

type args struct {
	Yaml string `arg:"required,positional"`
}

func (args) Version() string {
	return version
}

// YAML specification (basically)

type environment struct {
	Title string `yaml:"environment"`
	//TODO: implement other fields for the YAML
	//Run []Action
}

type environmentList map[string]environment

type environmentsConfig struct {
	SessionUser string   `yaml:"default_username"`
	Environments       environmentList `yaml:"environments"`
}

// Function to interpret the YAML file
func (um *environmentList) UnmarshalYAML(val *yaml.Node) error {
	var slice []environment
	if err := val.Decode(&slice); err != nil {
		return err
	}

	*um = make(environmentList)
	for _, env := range slice {
		(*um)[env.Title] = env
	}

	return nil
}


// TUI stuff

// The app will revolve around views that can have a rsult, send tea.Cmd's, and
// be able to `push` a new view or `pop` itself
type viewResult struct {
	view view
	cmd  tea.Cmd
	push view
	pop  bool
}

type view interface {
	View(width, height int) string
	Update(msg tea.Msg)     viewResult
}

type mainModel struct {
	width, height int
	viewStack     []view
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) View() string {
	if len(m.viewStack) == 0 {
		return ""
	}

	current := m.viewStack[len(m.viewStack) - 1]
	return current.View(m.width, m.height)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
	}

	result := m.viewStack[len(m.viewStack) - 1].Update(msg)

	if result.pop {
		m.viewStack = m.viewStack[:len(m.viewStack) - 1]
	}
	if result.push != nil {
    		m.viewStack = append(m.viewStack, result.push)
    	} else if !result.pop && result.push == nil {
		m.viewStack[len(m.viewStack) - 1] = result.view
	}

	if len(m.viewStack) == 0 {
		return m, tea.Quit
	}

	return m, nil
}

// The main and initial view
type menuView struct {
	items  []string
	cursor int
}

func (v menuView) View(width, height int) string {
	var content strings.Builder
	content.WriteString("available environments:\n\n")
	for _, item := range v.items {
    		if v.items[v.cursor] == item {
	    		content.WriteString("\u2192")
    		}
		content.WriteString(" ")
		content.WriteString(item)
		content.WriteString("\n")
	}
	content.WriteString("\nuse ctrl+q to quit.\n")

	return content.String()
}

func (v menuView) Update(msg tea.Msg) viewResult {
	switch msg := msg.(type) {
	case tea.KeyMsg:
    		switch msg.String() {
		case "ctrl+q":
	    		return viewResult{view: v, pop: true}
		case "up", "k":
			if v.cursor != 0 {
				v.cursor--
			}
		case "down", "j":
	    		if v.cursor != len(v.items) - 1 {
				v.cursor++
	    		}
		case "enter":
            		return viewResult{
				// I'm not sure if we even need to always return the view
				// while we're pushing another? TODO: test this
				view: v,
				// TODO: update the return when other functions are implemented
				push: errorView{message: v.items[v.cursor]},
	    		}
		}
	}
	return viewResult{view: v}
}

//Hopefully a functional "view" that spawns a shell
type interactiveView struct {
	command string
}

func (v interactiveView) View() {

}

func (v interactiveView) Update() {

}

// A view made to display something like an error message
type errorView struct {
	message string
}

func (v errorView) View(width, height int) string {
	content := "There's been an error:\n\n" + v.message
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func (v errorView) Update(msg tea.Msg) viewResult {
	// For now, any keypress will close the errorView
	if _, ok := msg.(tea.KeyMsg); ok {
		return viewResult{pop: true}
	}
	return viewResult{view: v}
}

// Init & Main

var version = "dev"
var environments_yaml []byte
var cfg environmentsConfig

func init() {
	var args args
	arg.MustParse(&args)

	var err error

	// Check if command-line specified YAML file is readable
	if _, err = os.Stat(args.Yaml); err != nil {
    		log.Fatal(err)
	}

	// Read the file
    	if environments_yaml, err = os.ReadFile(args.Yaml); err != nil {
        	log.Fatal(err)
    	}

	// Create environments config/list struct
	if err = yaml.Unmarshal(environments_yaml, &cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	m := mainModel{
		viewStack: []view{
			menuView{items: []string{"hotdog", "burger"},},
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
