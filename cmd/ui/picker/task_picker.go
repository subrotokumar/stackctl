package picker

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/internal/task"
)

var (
	pickTitleStyle    = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	pickCursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	pickSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	pickDescStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	pickHelpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

type pickerModel struct {
	names  []string
	descs  map[string]string
	width  int // width of the name column
	cursor int
	chosen string
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "ctrl+c", "esc", "q":
		return m, tea.Quit
	case "up", "k":
		m.cursor = (m.cursor - 1 + len(m.names)) % len(m.names)
	case "down", "j":
		m.cursor = (m.cursor + 1) % len(m.names)
	case "home", "g":
		m.cursor = 0
	case "end", "G":
		m.cursor = len(m.names) - 1
	case "enter":
		m.chosen = m.names[m.cursor]
		return m, tea.Quit
	}
	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder
	b.WriteString(core.PurpleStyle.Render("Tasks") + "\n")
	b.WriteString(core.GreyStyle.Render("Select a task to run") + "\n")
	for i, n := range m.names {
		name := fmt.Sprintf("%-*s", m.width, n)
		desc := core.GreyStyle.Render(m.descs[n])
		if i == m.cursor {
			b.WriteString(pickCursorStyle.Render("❯ ") + (name) + "  " + desc + "\n")
		} else {
			b.WriteString("  " + name + "  " + desc + "\n")
		}
	}
	b.WriteString(pickHelpStyle.Render("↑/↓ navigate • enter run • q quit") + "\n")
	return b.String()
}

// PickTask shows an interactive selector and returns the chosen task name.
// It returns "" with a nil error if the user cancels.
func PickTask(tf *task.Taskfile) (string, error) {
	if len(tf.Tasks) == 0 {
		return "", fmt.Errorf("no tasks defined in task file")
	}
	if st, err := os.Stdin.Stat(); err != nil || st.Mode()&os.ModeCharDevice == 0 {
		return "", fmt.Errorf("no task given and stdin is not a terminal; use: stackctl run <task>")
	}

	m := pickerModel{descs: map[string]string{}}
	for n, t := range tf.Tasks {
		m.names = append(m.names, n)
		m.descs[n] = t.Desc
		if len(n) > m.width {
			m.width = len(n)
		}
	}
	sort.Strings(m.names)

	final, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", err
	}
	return final.(pickerModel).chosen, nil
}
