package panels

import (
	"fmt"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/alerts"
	"github.com/Dicklesworthstone/ntm/internal/tui/layout"
	"github.com/Dicklesworthstone/ntm/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AlertsPanel struct {
	width   int
	height  int
	focused bool
	alerts  []alerts.Alert
}

func NewAlertsPanel() *AlertsPanel {
	return &AlertsPanel{}
}

func (m *AlertsPanel) SetData(alerts []alerts.Alert) {
	m.alerts = alerts
}

func (m *AlertsPanel) Init() tea.Cmd {
	return nil
}

func (m *AlertsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AlertsPanel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *AlertsPanel) Focus() {
	m.focused = true
}

func (m *AlertsPanel) Blur() {
	m.focused = false
}

func (m *AlertsPanel) View() string {
	t := theme.Current()
	
	if m.width <= 0 {
		return ""
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Text).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(t.Surface1).
		Width(m.width).
		Padding(0, 1).
		Render("Active Alerts")

	var content strings.Builder
	content.WriteString(header + "\n")

	if len(m.alerts) == 0 {
		content.WriteString("\n  " + lipgloss.NewStyle().Foreground(t.Green).Render("No active alerts") + "\n")
		return content.String()
	}

	for i, a := range m.alerts {
		if i >= 10 { break } // limit display
		
		color := t.Blue
		icon := "ℹ"
		switch a.Severity {
		case alerts.SeverityCritical:
			color = t.Red
			icon = "✗"
		case alerts.SeverityWarning:
			color = t.Yellow
			icon = "⚠"
		}
		
		msg := layout.TruncateRunes(a.Message, m.width-6, "...")
		line := fmt.Sprintf("  %s %s", icon, msg)
		content.WriteString(lipgloss.NewStyle().Foreground(color).Render(line) + "\n")
	}

	return content.String()
}
