package tui

import "github.com/charmbracelet/lipgloss"

type Styler struct {
	titleStyle       lipgloss.Style
	highlightStyle   lipgloss.Style
	cursorStyle      lipgloss.Style
	defaultStyle     lipgloss.Style
	selectedDNSStyle lipgloss.Style
	infoStyle        lipgloss.Style
	instructionStyle lipgloss.Style
}

func NewDefaultStyler() *Styler {
	return &Styler{
		titleStyle:       lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0e69c9")) /*.Padding(0, 1)*/,
		highlightStyle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0e69c9")),
		cursorStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0e69c9")),
		defaultStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
		selectedDNSStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#ADFF2F")),
		infoStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#00CED1")),
		instructionStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#B0C4DE")),
	}
}
