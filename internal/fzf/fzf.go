package fzf

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/structs"
)

type model struct {
	aliases []structs.DirAlias
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
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Which project do you want to open?\n\n"

	for i, a := range m.aliases {
		if i == 0 {
			s += fmt.Sprintf("> %s  %s <\n\n", a.Alias, a.Path)
			continue
		}

		s += fmt.Sprintf("  %s  %s\n", a.Alias, a.Path)

		if i >= 9 {
			break
		}
	}

	s += "\nPress q to quit.\n"
	return s
}

func StartFzf(configPath string) *tea.Program {
	return tea.NewProgram(initialModel(configPath))
}
