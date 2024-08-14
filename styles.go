package gologging

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	TimeStyle  lipgloss.Style
	TagStyle   lipgloss.Style
	PkgStyle   lipgloss.Style
	MsgStyle   lipgloss.Style
	BlockStyle lipgloss.Style
	FileStyle  lipgloss.Style
	IconStyle  lipgloss.Style
}

func defaultStyles(r *lipgloss.Renderer) *Styles {
	return &Styles{
		TimeStyle:  lipgloss.NewStyle().MarginRight(2).Renderer(r),
		TagStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#62c6b7")).Width(10).Renderer(r),
		PkgStyle:   lipgloss.NewStyle().Width(24).Renderer(r),
		MsgStyle:   lipgloss.NewStyle().Renderer(r),
		BlockStyle: lipgloss.NewStyle().MarginLeft(1).Renderer(r),
		FileStyle:  lipgloss.NewStyle().Margin(0, 1).Foreground(lipgloss.Color("#61afe1")).Renderer(r),
		IconStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Padding(0, 1).MarginRight(1).Renderer(r),
	}
}
