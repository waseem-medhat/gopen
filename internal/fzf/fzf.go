// Package fzf contains types and logic for the interactive fuzzy finder part
// of Gopen.
package fzf

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	l "github.com/charmbracelet/lipgloss"
	"github.com/wipdev-tech/gopen/internal/config"
)

var styles = struct {
	selected l.Style
	rest     l.Style
	cursor   l.Style
	window   l.Style
	question l.Style
	logo     l.Style
}{
	logo:     l.NewStyle().Foreground(l.Color("56")),
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
			m.done = true
			m.Selected = ""
			return m, tea.Quit

		case "ctrl+w":
			m.searchStr = ""

		case "up", "ctrl+p":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}

		case "down", "ctrl+n":
			if m.selectedIdx < 9 && m.selectedIdx < len(m.Config.DirAliases)-1 {
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

	logo := `
   _____                            
  / ____|                           
 | |  __   ___   _ __    ___  _ __  
 | | |_ | / _ \ | '_ \  / _ \| '_ \ 
 | |__| || (_) || |_) ||  __/| | | |
  \_____| \___/ | .__/  \___||_| |_|
                | |                 
                |_|                 
`

	maxLenAlias := 0
	maxLenPath := 0
	for _, a := range m.Config.DirAliases {
		if len(a.Alias) > maxLenAlias {
			maxLenAlias = len(a.Alias)
		}
		if len(a.Path) > maxLenPath {
			maxLenPath = len(a.Path)
		}
	}

	fmtStr := fmt.Sprintf("%%-%ds", maxLenAlias+maxLenPath+6)
	s := styles.question.Render(
		fmt.Sprintf(fmtStr, "Which project do you want to open?"),
	)
	s += fmt.Sprintf("\n\n> %s", m.searchStr)
	s += styles.cursor.Render("█")
	s += "\n\n"

	fmtStr = fmt.Sprintf("  %%-%ds  %%-%ds ", maxLenAlias, maxLenPath+1)
	for i, a := range m.results {
		if i == m.selectedIdx {
			s += styles.selected.Render(fmt.Sprintf(fmtStr, a.Alias, a.Path))
			s += "\n"
			continue
		}

		s += styles.rest.Render(fmt.Sprintf(fmtStr, a.Alias, a.Path))
		s += "\n"

		if i > 9 {
			break
		}
	}

	if m.helpShown {
		s += "\n?         hide key bindings"
		s += "\nctrl+n/↓  move selection down"
		s += "\nctrl+p/↑  move selection up"
		s += "\nctrl+w    clear search string"
		s += "\nctrl+c    quit"
	} else {
		s += "\n?         show key bindings"
		s += "\nctrl+c    quit"
	}
	return styles.logo.Render(logo) + "\n" + styles.window.Render(s) + "\n\n"
}

func searchAliases(aliases []config.DirAlias, searchStr string) []config.DirAlias {
	newResults := []config.DirAlias{}
	for _, a := range aliases {
		if strings.Contains(a.Alias, searchStr) || strings.Contains(a.Path, searchStr) {
			newResults = append(newResults, a)
		}

		if len(newResults) >= 10 {
			break
		}
	}
	return newResults
}

func initialModel(configPath string) Model {
	cfg, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}
	return Model{
		Config:  cfg,
		results: cfg.DirAliases[0:10],
	}
}

// StartFzf is the entry point for the fuzzy finder which spawns the bubbletea
// program.
func StartFzf(configPath string) *tea.Program {
	return tea.NewProgram(initialModel(configPath))
}
