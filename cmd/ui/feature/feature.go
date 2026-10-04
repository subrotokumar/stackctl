package feature

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/v4/cmd/core"
	"github.com/subrotokumar/stackctl/v4/internal/init/micronaut"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const selectorIndicator = "✅"

type item struct {
	f        micronaut.Feature
	selected bool
}

func (i item) Title() string {
	title := fmt.Sprintf("%s %s", i.f.Title, core.BlueStyle.Render("["+i.f.Name+"]"))
	if i.selected {
		title = selectorIndicator + " " + title
	}
	if i.f.Preview {
		title += " " + core.PreviewCode.Render("Preview")
	}
	if i.f.Community {
		title += " " + core.GreyStyle.Render("(community)")
	}
	return title
}

func (i item) Description() string {
	if i.f.Category == "" {
		return i.f.Description
	}
	return core.GreyStyle.Render(i.f.Category) + "  " + i.f.Description
}

func (i item) FilterValue() string { return i.f.Title + " " + i.f.Name + " " + i.f.Category }

type model struct {
	list      list.Model
	features  []micronaut.Feature
	selected  core.Set[string]
	back      bool
	cancelled bool
}

func New(features []micronaut.Feature) model {
	m := model{features: features, selected: make(core.Set[string])}

	items := make([]list.Item, 0, len(features))
	for _, f := range features {
		items = append(items, item{f: f})
	}

	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Features"
	m.list.DisableQuitKeybindings() // enter / esc / ctrl+c are handled here
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "select")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		}
	}
	return m
}

// WithSelected pre-selects the given feature names, so going back to this
// step shows the earlier picks.
func (m model) WithSelected(names []string) model {
	for _, n := range names {
		m.selected.Add(n)
	}
	for i, f := range m.features {
		if m.selected.Has(f.Name) {
			m.list.SetItem(i, item{f: f, selected: true})
		}
	}
	return m
}

func (m *model) Init() tea.Cmd { return nil }

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
				f := m.features[index]
				if m.selected.Has(f.Name) {
					m.selected.Remove(f.Name)
				} else {
					m.selected.Add(f.Name)
				}
				m.list.SetItem(index, item{f: f, selected: m.selected.Has(f.Name)})
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

func (m *model) View() string { return docStyle.Render(m.list.View()) }

// RunWithBack returns the selected feature names, and back=true if the user
// pressed esc (with no filter active). Ctrl+C exits the program.
func (m model) RunWithBack() ([]string, bool) {
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	if m.cancelled {
		os.Exit(0)
	}
	return m.selected.ToSlice(), m.back
}
