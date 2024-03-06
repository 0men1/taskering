package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var Style = lipgloss.NewStyle().Bold(true)

var CatStyle = lipgloss.NewStyle().
			Underline(true).
			Border(lipgloss.NormalBorder(), true, true, true, true)

