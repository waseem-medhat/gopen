// Package tui contains types and logic for the interactive TUI part of Gopen.
package tui

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	l "github.com/charmbracelet/lipgloss"
	"github.com/waseem-medhat/gopen/internal/config"
)

var styles = struct {
	selected l.Style
	rest     l.Style
	cursor   l.Style
	window   l.Style
	question l.Style
	logo     l.Style
}{
	logo:     l.NewStyle().Foreground(l.Color("57")),
	question: l.NewStyle().Bold(true),
	rest:     l.NewStyle().Faint(true),
	cursor:   l.NewStyle().Blink(true),
	window: l.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Border(l.RoundedBorder()),
	selected: l.NewStyle().
		Bold(true).
		Foreground(l.Color("255")).
		Background(l.Color("56")),
}

// Model implements the tea.Model interface to be used as the model part of the
// bubbletea program and includes fields that hold the program state.
//
// Note that the fields `Config` and `Selected` are exported because the are
// used by the main package.
type Model struct {
	Config      config.C
	Selected    string
	searchStr   string
	results     []config.DirAlias
	selectedIdx int
	helpShown   bool
	done        bool
}

// Init is one of the tea.Model interface methods but not used by the TUI.
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
			m.done = true
			m.Selected = ""
			return m, tea.Quit

		case "ctrl+w":
			m.searchStr = ""
			m.results = searchAliases(m.Config.DirAliases, m.searchStr)

		case "up", "ctrl+p":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "ctrl+n":
			if m.selectedIdx < len(m.results)-1 {
				m.selectedIdx++
			}

		case "enter":
			m.done = true
			return m, tea.Quit

		case "backspace":
			if len(m.searchStr) >= 1 {
				m.searchStr = m.searchStr[:len(m.searchStr)-1]
			}
			m.results = searchAliases(m.Config.DirAliases, m.searchStr)
			m.selectedIdx = 0

		case "?":
			m.helpShown = !m.helpShown

		default:
			if len(msg.String()) == 1 {
				m.searchStr += msg.String()
				m.results = searchAliases(m.Config.DirAliases, m.searchStr)
				m.selectedIdx = 0
			}
		}
	}

	if len(m.results) > 0 {
		m.Selected = m.results[m.selectedIdx].Alias
	}
	return m, nil
}

// View is one of the tea.Model interface methods. It includes the rendering logic.
func (m Model) View() string {
	if m.done {
		return ""
	}

	maxAliasW, maxPathW, maxW := calcMaxWidths(m.Config.DirAliases)
	logo := styles.logo.Render(gopenLogo)
	question := styles.question.Render(
		alignQuestion("Which project do you want to open?", maxW),
	)
	promptLine := fmt.Sprintf("> %s", m.searchStr) + styles.cursor.Render("█")

	results := ""
	for i, a := range m.results {
		if i == m.selectedIdx {
			results += styles.selected.Render(
				alignResult(a.Alias, a.Path, maxAliasW, maxPathW),
			)
		} else {
			results += styles.rest.Render(
				alignResult(a.Alias, a.Path, maxAliasW, maxPathW),
			)
		}

		results += "\n"
		if i >= 5 {
			break
		}
	}

	window := styles.window.Render(question + "\n\n" + promptLine + "\n\n" + results)

	help := ""
	if m.helpShown {
		help = fullHelp
	} else {
		help = shortHelp
	}

	return logo + "\n" + window + help + "\n\n"
}

func searchAliases(aliases []config.DirAlias, searchStr string) []config.DirAlias {
	newResults := []config.DirAlias{}
	for _, a := range aliases {
		if strings.Contains(a.Alias, searchStr) || strings.Contains(a.Path, searchStr) {
			newResults = append(newResults, a)
		}

		if len(newResults) >= 5 {
			break
		}
	}
	return newResults
}

func initialModel(cfg config.C) Model {
	results := cfg.DirAliases
	slices.Reverse(results)
	if len(cfg.DirAliases) > 5 {
		results = cfg.DirAliases[0:5]
	}

	return Model{
		Config:  cfg,
		results: results,
	}
}

// StartTUI is the entry point for the interactive TUI which spawns the
// bubbletea program.
func StartTUI(cfg config.C) *tea.Program {
	return tea.NewProgram(initialModel(cfg))
}
