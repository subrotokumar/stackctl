package extension

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/v4/cmd/core"
	"github.com/subrotokumar/stackctl/v4/internal/init/quarkus"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var selected core.Set[string] = make(core.Set[string])

const SELECTOR_INDICATOR string = "✅"

type ItemType = item

type item struct {
	*quarkus.Extension
	title, desc string
}

func NewItem(param quarkus.Extension) item {
	return item{
		title:     param.Name,
		desc:      param.Description,
		Extension: &param,
	}
}

func (i item) Title() string {
	id := core.BlueStyle.Render(fmt.Sprintf("[%s]", strings.Split(i.ID, ":")[1]))

	isSelected := ""
	platform := ""

	if i.Platform {
		platform = core.SelectedStyle.Render(" ※ ")
	}

	if selected.Has(i.ID) {
		isSelected = SELECTOR_INDICATOR + " "
	}

	tag := ""
	for _, val := range i.Tags {
		switch val {
		case "with:starter-code":
			tag = fmt.Sprintf("%s %s ", tag, core.StartedCode.Render("(Starter Code)"))
		case "status:stable":
			tag = tag + ""
		case "status:deprecated":
			tag = fmt.Sprintf("%s %s ", tag, core.DeprecatedCode.Render("Deprecated"))
		case "status:preview":
			tag = fmt.Sprintf("%s %s ", tag, core.PreviewCode.Render("Preview"))
		case "status:experimental":
			tag = fmt.Sprintf("%s %s ", tag, core.ExperimentalCode.Render("Experimental"))
		default:
		}
	}

	return fmt.Sprintf("%s%s %s%s%s", isSelected, i.Name, id, platform, tag)
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.Name
}
func (i item) Id() string { return i.ID }

type model struct {
	list list.Model

	options   []quarkus.Extension
	back      bool
	cancelled bool
}

func New(options []quarkus.Extension) model {
	m := model{options: options}
	inputOptions := []list.Item{}

	for _, val := range m.options {
		inputOptions = append(inputOptions, NewItem(val))
	}

	m.list = list.New(inputOptions, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Dependencies"
	m.list.DisableQuitKeybindings() // we handle enter / esc / ctrl+c ourselves
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "select")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		}
	}
	selected = make(core.Set[string])
	return m
}

// WithSelected pre-selects the given extension IDs, so going back
// to this step shows the earlier picks.
func (m model) WithSelected(ids []string) model {
	for _, id := range ids {
		selected.Add(id)
	}
	for i := range m.options {
		if selected.Has(m.options[i].ID) {
			m.options[i].Selected = true
			m.list.SetItem(i, NewItem(m.options[i]))
		}
	}
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.cancelled = true
			return m, tea.Quit
		}

		// While the user is typing a filter, enter/esc/space belong to the list.
		if m.list.FilterState() != list.Filtering {
			switch msg.Type {
			case tea.KeyEnter:
				return m, tea.Quit

			case tea.KeyEsc:
				// With a filter applied, esc falls through and clears it.
				if m.list.FilterState() == list.Unfiltered {
					m.back = true
					return m, tea.Quit
				}

			case tea.KeySpace:
				index := m.list.GlobalIndex()
				m.options[index].Selected = !m.options[index].Selected
				if m.options[index].Selected {
					selected.Add(m.options[index].ID)
				} else {
					selected.Remove(m.options[index].ID)
				}
				m.list.SetItem(index, NewItem(m.options[index]))
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return docStyle.Render(m.list.View())
}

func (m model) Run() []string {
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return selected.ToSlice()
}

// RunWithBack returns the selected extension IDs, and back=true if the
// user pressed esc (with no filter active). Ctrl+C exits the program.
func (m model) RunWithBack() ([]string, bool) {
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if m.cancelled {
		os.Exit(0)
	}
	return selected.ToSlice(), m.back
}
