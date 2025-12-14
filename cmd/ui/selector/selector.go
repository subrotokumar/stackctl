package selector

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrotokumar/stackctl/cmd/core"
)

type model struct {
	title    string
	choice   int
	option   []string
	quitting bool
}

func New(title string, option []string) model {
	return model{
		title:  title,
		option: option,
	}
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			os.Exit(0)
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

	if m.quitting {
		return ""
	}

	return s
}

func (m model) Run() string {
	_, _ = tea.NewProgram(&m).Run()
	return m.option[m.choice]
}
