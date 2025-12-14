package extension

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/subrotokumar/stackctl/cmd/core"
	"github.com/subrotokumar/stackctl/internal/quarkus"
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
	id := core.BlueStyle.Render(fmt.Sprintf("[%s]", core.BlueStyle.Render(strings.Split(i.ID, ":")[1])))

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
			tag = fmt.Sprintf("%s %s ", tag, core.StartedCode.Render("Starter Code"))
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

	options []quarkus.Extension
}

func New(options []quarkus.Extension) model {
	m := model{options: options}
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
