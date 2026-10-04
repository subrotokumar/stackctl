package selector

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrotokumar/stackctl/v4/cmd/core"
)

type model struct {
	title    string
	choice   int
	option   []string
	quitting bool
	back     bool
}

func New(title string, option []string) model {
	return model{
		title:  title,
		option: option,
	}
}

// WithDefault pre-selects the option equal to value (if present).
// Used so that going back shows the previously chosen answer.
func (m model) WithDefault(value string) model {
	for i, o := range m.option {
		if o == value {
			m.choice = i
			break
		}
	}
	return m
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			os.Exit(0)
			return m, tea.Quit

		case "esc", "ctrl+b":
			m.quitting = true
			m.back = true
			return m, tea.Quit

		case "up":
			if m.choice > 0 {
				m.choice--
			}

		case "down":
			if m.choice < len(m.option)-1 {
				m.choice++
			}

		case "enter":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}

	s := "\n" + core.QuestionStyle.Render(m.title) + ":\n"

	for i, opt := range m.option {
		cursor := "[ ]"
		style := core.UnSelectedStyle

		if i == m.choice {
			cursor = "[X]"
			style = core.SelectedStyle
		}

		s += style.Render(cursor+" "+opt) + "\n"
	}

	s += "\n" + core.UnSelectedStyle.Render("↑/↓: move • enter: confirm • esc: back • ctrl+c: quit") + "\n"
	return s
}

func (m model) Run() string {
	_, _ = tea.NewProgram(&m).Run()
	return m.option[m.choice]
}

// RunWithBack returns the chosen option, and back=true if the user pressed esc/ctrl+b.
func (m model) RunWithBack() (string, bool) {
	_, _ = tea.NewProgram(&m).Run()
	return m.option[m.choice], m.back
}
