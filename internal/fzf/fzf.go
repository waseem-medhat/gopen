package fzf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wipdev-tech/gopen/internal/config"
)

type FzfModel struct {
	Config config.C
	Prompt string
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

func initialModel(configPath string) FzfModel {
	cfg, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}
	return FzfModel{
		Config: cfg,
	}
}

func (m FzfModel) Init() tea.Cmd {
	return nil
}

func (m FzfModel) cmdGopen() tea.Msg {
	return tea.QuitMsg{}
}

func (m FzfModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+w":
			m.Prompt = ""
		case "enter":
			return m, m.cmdGopen
		case "backspace":
			if len(m.Prompt) >= 1 {
				m.Prompt = m.Prompt[:len(m.Prompt)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.Prompt += msg.String()
			}
		}
	}

	return m, nil
}

func (m FzfModel) View() string {
	s := fmt.Sprintf("Which project do you want to open?\n> %s", m.Prompt)
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

func StartFzf(configPath string) *tea.Program {
	return tea.NewProgram(initialModel(configPath))
}
