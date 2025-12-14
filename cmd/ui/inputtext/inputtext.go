package inputtext

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrotokumar/stackctl/cmd/core"
)

var hidePlaceHolder = true

type (
	errMsg error
)

type model struct {
	title        string
	defaultValue string
	textInput    textinput.Model
	err          error
	quitting     bool
}

func New(title string, defaultValue string) model {
	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 200

	return model{
		title:        title,
		defaultValue: defaultValue,
		textInput:    ti,
		err:          nil,
		quitting:     false,
	}
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyEnter:
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.quitting {
		return ""
	}
	if hidePlaceHolder {
		return fmt.Sprintf(
			"\n%s :\n%s",
			core.QuestionStyle.Render(m.title),
			m.textInput.View(),
		) + "\n"
	}
	return fmt.Sprintf(
		"\n%s (%s):\n%s",
		core.QuestionStyle.Render(m.title),
		core.GreyStyle.Render(m.defaultValue),
		m.textInput.View(),
	) + "\n"
}

func (m model) Run() string {
	_, _ = tea.NewProgram(&m).Run()
	if val := m.textInput.Value(); val == "" {
		return m.defaultValue
	} else {
		return val
	}
}
