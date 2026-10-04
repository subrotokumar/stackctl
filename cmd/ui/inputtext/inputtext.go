package inputtext

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/subrotokumar/stackctl/v4/cmd/core"
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
	back         bool
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
		case tea.KeyCtrlC:
			os.Exit(0)
			return m, tea.Quit
		case tea.KeyEsc:
			m.quitting = true
			m.back = true
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
	hint := core.GreyStyle.Render("enter: confirm • esc: back • ctrl+c: quit")
	if hidePlaceHolder {
		return fmt.Sprintf(
			"\n%s :\n%s",
			core.QuestionStyle.Render(m.title),
			m.textInput.View(),
		) + "\n\n" + hint + "\n"
	}
	return fmt.Sprintf(
		"\n%s (%s):\n%s",
		core.QuestionStyle.Render(m.title),
		core.GreyStyle.Render(m.defaultValue),
		m.textInput.View(),
	) + "\n\n" + hint + "\n"
}

func (m model) value() string {
	if val := m.textInput.Value(); val != "" {
		return val
	}
	return m.defaultValue
}

func (m model) Run() string {
	_, _ = tea.NewProgram(&m).Run()
	return m.value()
}

// RunWithBack returns the entered text (or the default if empty),
// and back=true if the user pressed esc.
func (m model) RunWithBack() (string, bool) {
	_, _ = tea.NewProgram(&m).Run()
	return m.value(), m.back
}
