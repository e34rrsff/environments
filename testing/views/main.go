package main

import (
    "os"
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type View interface {
    View(width, height int) string
    Update(msg tea.Msg) ViewResult
}

type ViewResult struct {
    View View
    cmd tea.Cmd
    Push View
    Pop bool
}

// AVAILABLE APP VIEWS
// |SelectionView
type SelectionView struct {
    message string
}

func (v SelectionView) View(width, height int) string {
    content := "Selection:\n\n" + v.message
    return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func (v SelectionView) Update(msg tea.Msg) ViewResult {
    if _, ok := msg.(tea.KeyMsg); ok {
        return ViewResult{Pop: true}
    }
    return ViewResult{View: v}
}
//
// |MenuView
type MenuView struct {
    items []string
    cursor int
}

func (v MenuView) View(width, height int) string {
    var content strings.Builder

    content.WriteString("Available Exams:\n\n") 
    for _, item := range v.items {
        if v.items[v.cursor] == item {
            content.WriteString("\u2192")
        }
        content.WriteString(" ")
        content.WriteString(item) 
        content.WriteString("\n")
    }

    return content.String()
}

func (v MenuView) Update(msg tea.Msg) ViewResult {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {

        case "ctrl+q":
            return ViewResult{View: v, Pop: true}

        case "up", "k":
            if v.cursor != 0 {
                v.cursor--
            }

        case "down", "j":
            if v.cursor != len(v.items) - 1 {
                v.cursor++
            }

        case "enter":
            return ViewResult{View: v, Push: SelectionView{message: v.items[v.cursor]}, /*Pop: true*/} 
        }
    }

    return ViewResult{View: v}
}

type model struct {
    width, height int
    viewStack     []View
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) View() string {
    if len(m.viewStack) == 0 {
        return ""
    }

    current := m.viewStack[len(m.viewStack) - 1]
    return current.View(m.width, m.height)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil
    }

    current := m.viewStack[len(m.viewStack) - 1]
    result := current.Update(msg)

    if result.Pop {
        m.viewStack = m.viewStack[:len(m.viewStack) - 1]
    }
    if result.Push != nil {
        m.viewStack = append(m.viewStack, result.Push)
    } else if !result.Pop && result.Push == nil {
        m.viewStack[len(m.viewStack) - 1] = result.View
    }

    if len(m.viewStack) == 0 {
        return m, tea.Quit
    }

    return m, nil
}

func main() {
    m := model{
        viewStack: []View{
            MenuView{items: []string{"hotdog", "burger"},},
        },
    }

    p := tea.NewProgram(m, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
