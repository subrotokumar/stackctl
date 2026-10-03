package listview

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/internal/init/spring"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var selected core.Set[string] = make(core.Set[string])

const SELECTOR_INDICATOR string = "✅"

type ItemType = item

type item struct {
	*spring.DependencyDetail
	title, desc string
}

func NewItem(param spring.DependencyDetail) item {
	return item{
		title:            param.Name,
		desc:             param.Description,
		DependencyDetail: &param,
	}
}

func (i item) Title() string {
	tag := "(" + core.GreyStyle.Render(i.Tag) + ")"
	if selected.Has(i.DependencyDetail.ID) {
		return fmt.Sprintf("%s %s  %s", SELECTOR_INDICATOR, i.FilterValue(), core.GreyStyle.Render(tag))
	}
	return fmt.Sprintf("%s  %s", i.title, core.GreyStyle.Render(tag))
}

func (i item) Description() string {
	if i.DependencyDetail.VersionRange != nil {
		versionRange := *i.DependencyDetail.VersionRange
		return fmt.Sprintf("%s %s", i.desc, core.RedStyle.Render(versionRange))
	}
	return i.desc
}

func (i item) FilterValue() string { return i.title }
func (i item) Id() string          { return i.ID }

type model struct {
	list list.Model

	options   []spring.DependencyDetail
	back      bool
	cancelled bool
}

func New(options []spring.DependencyGroup) model {
	detailList := []spring.DependencyDetail{}

	for _, dependencyGroup := range options {
		tag := dependencyGroup.Name
		for _, detail := range dependencyGroup.Values {
			detail.Tag = tag
			detailList = append(detailList, detail)
		}
	}
	m := model{options: detailList}
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

// WithSelected pre-selects the given dependency IDs, so going back
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

// RunWithBack returns the selected dependency IDs, and back=true if the
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
