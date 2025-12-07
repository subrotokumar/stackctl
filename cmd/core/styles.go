package core

import "github.com/charmbracelet/lipgloss"

var (
	LogoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	QuestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01bcfaff")).Bold(true)

	SelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#f476ffff")).Bold(true)
	UnSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#b5b5b5ff")).Bold(true)
	GreyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#7a7a7aff")).Bold(true)

	RedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3a3aff")).Bold(false)

	TipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	EndingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)
