package core

import "github.com/charmbracelet/lipgloss"

var (
	GreenStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#01fa65")). // text
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#01fa65")). // border
			Bold(true).
			Padding(0, 2)

	LogoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	QuestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01bcfaff")).Bold(true)
	BlueStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#73bffcff")).Bold(true)

	StartedCode = lipgloss.NewStyle().
			Margin(0).
			Padding(0, 1).
			Background(lipgloss.Color("#fcfc7393")).
			Foreground(lipgloss.Color("#ffffffff")).
			Bold(true)

	DeprecatedCode = lipgloss.NewStyle().
			Margin(0).
			Padding(0, 1).
			Background(lipgloss.Color("#757575ff")).
			Foreground(lipgloss.Color("#ffffffff")).
			Bold(true)

	PreviewCode = lipgloss.NewStyle().
			Margin(0).
			Padding(0, 1).
			Background(lipgloss.Color("#2a83ffff")).
			Foreground(lipgloss.Color("#ffffffff")).
			Bold(true)

	ExperimentalCode = lipgloss.NewStyle().
				Margin(0).
				Padding(0, 1).
				Background(lipgloss.Color("#ff5219ff")).
				Foreground(lipgloss.Color("#ffffffff")).
				Bold(true)

	SelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f476ffff")).Bold(true)
	UnSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#b5b5b5ff")).Bold(true)
	GreyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#7a7a7aff")).Bold(true)

	RedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3a3aff")).Bold(false)

	TipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	EndingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)
