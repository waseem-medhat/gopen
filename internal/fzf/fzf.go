// Package fzf contains types and logic for the interactive fuzzy finder part
// of Gopen.
package fzf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wipdev-tech/gopen/internal/config"
)

var styles = struct {
	first  lipgloss.Style
	rest   lipgloss.Style
	prompt lipgloss.Style
}{
	first:  lipgloss.NewStyle().Foreground(lipgloss.Color("37")),
	rest:   lipgloss.NewStyle().Faint(true),
	prompt: lipgloss.NewStyle().Blink(true),
}

func initialModel(configPath string) Model {
	cfg, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}
	return Model{
		Config: cfg,
	}
}

// Model implements the tea.Model interface to be used as the model part of the
// bubbletea program, but includes fields that hold the program state, namely
// the config data and the search string
type Model struct {
	Config    config.C
	SearchStr string
}

// Init is one of the tea.Model interface methods but not used by the fuzzy
// finder.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is one of the tea.Model interface methods. It triggers updates to the
// model and its state on keypresses.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+w":
			m.SearchStr = ""
		case "enter":
			return m, tea.Quit
		case "backspace":
			if len(m.SearchStr) >= 1 {
				m.SearchStr = m.SearchStr[:len(m.SearchStr)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.SearchStr += msg.String()
			}
		}
	}

	return m, nil
}

// View is one of the tea.Model interface methods. It includes the rendering logic.
func (m Model) View() string {
	s := fmt.Sprintf("Which project do you want to open?\n> %s", m.SearchStr)
	s += styles.prompt.Render("|")
	s += "\n\n"

	maxLen := 0
	for _, a := range m.Config.DirAliases {
		if len(a.Alias) > maxLen {
			maxLen = len(a.Alias)
		}
	}

	for i, a := range m.Config.DirAliases {
		if i == 0 {
			fmtStr := fmt.Sprintf("[ %%-%ds  %%s ]", maxLen)
			s += styles.first.Render(fmt.Sprintf(fmtStr, a.Alias, a.Path))
			s += "\n"
			continue
		}

		fmtStr := fmt.Sprintf("  %%-%ds  %%s", maxLen)
		s += styles.rest.Render(fmt.Sprintf(fmtStr, a.Alias, a.Path))
		s += "\n"

		if i >= 9 {
			break
		}
	}

	s += "\nctrl+w: clear word"
	s += "\nctrl+c: quit\n"
	return s
}

// StartFzf is the entry point for the fuzzy finder which spawns the bubbletea
// program.
func StartFzf(configPath string) *tea.Program {
	return tea.NewProgram(initialModel(configPath))
}
