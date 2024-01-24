package fzf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/structs"
)

type model struct {
	aliases []structs.DirAlias
	prompt  string
}

var styles = struct {
	first  lipgloss.Style
	rest   lipgloss.Style
	prompt lipgloss.Style
}{
	first:  lipgloss.NewStyle().Foreground(lipgloss.Color("37")),
	rest:   lipgloss.NewStyle().Faint(true),
	prompt: lipgloss.NewStyle().Blink(true),
}

func initialModel(configPath string) model {
	cfg, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}
	return model{
		aliases: cfg.DirAliases,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+w":
			m.prompt = ""
		case "backspace":
			if len(m.prompt) >= 1 {
				m.prompt = m.prompt[:len(m.prompt)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.prompt += msg.String()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Which project do you want to open?\n> %s", m.prompt)
	s += styles.prompt.Render("|")
	s += "\n\n"

	maxLen := 0
	for _, a := range m.aliases {
		if len(a.Alias) > maxLen {
			maxLen = len(a.Alias)
		}
	}

	for i, a := range m.aliases {
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

func StartFzf(configPath string) *tea.Program {
	return tea.NewProgram(initialModel(configPath))
}
