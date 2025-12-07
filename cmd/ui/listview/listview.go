package listview

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/springx/cmd/core"
	"github.com/subrotokumar/springx/internal/spring"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

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
	if i.DependencyDetail.Selected {
		return fmt.Sprintf("%s  %s", core.SelectedStyle.Render(i.title), core.GreyStyle.Render(tag))
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

type model struct {
	list     list.Model
	selected core.Set[spring.DependencyDetail]
	options  []spring.DependencyDetail
}

func (m *model) UpdateList() {
	inputOptions := []list.Item{}

	sort.Slice(m.options, func(i, j int) bool {
		return !m.options[i].Selected && m.options[j].Selected
	})

	for _, val := range m.options {
		inputOptions = append(inputOptions, NewItem(val))
	}

	m.list = list.New(inputOptions, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Dependencies"
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
	m.UpdateList()
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeySpace:
			index := m.list.Index()
			m.options[index].Selected = !m.options[index].Selected
			m.list.SetItem(index, NewItem(m.options[index]))
			return m, nil
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

func (m model) Run() {
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
