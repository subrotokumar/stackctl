package listview

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/springx/cmd/core"
	"github.com/subrotokumar/springx/internal/spring"
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

	options []spring.DependencyDetail
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
	selected = make(core.Set[string])
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEnter, tea.KeyEsc:
			return m, tea.Quit
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
