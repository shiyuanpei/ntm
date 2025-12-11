package panels

import (
	"fmt"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/bv"
	"github.com/Dicklesworthstone/ntm/internal/tui/layout"
	"github.com/Dicklesworthstone/ntm/internal/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BeadsPanel struct {
	width   int
	height  int
	focused bool
	summary bv.BeadsSummary
	ready   []bv.BeadPreview
}

func NewBeadsPanel() *BeadsPanel {
	return &BeadsPanel{}
}

func (m *BeadsPanel) SetData(summary bv.BeadsSummary, ready []bv.BeadPreview) {
	m.summary = summary
	m.ready = ready
}

func (m *BeadsPanel) Init() tea.Cmd {
	return nil
}

func (m *BeadsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *BeadsPanel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *BeadsPanel) Focus() {
	m.focused = true
}

func (m *BeadsPanel) Blur() {
	m.focused = false
}

func (m *BeadsPanel) View() string {
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
		Render("Beads Pipeline")

	var content strings.Builder
	content.WriteString(header + "\n")

	// Stats row
	stats := fmt.Sprintf("Ready: %d  In Progress: %d  Blocked: %d",
		m.summary.Ready, m.summary.InProgress, m.summary.Blocked)
	statsStyled := lipgloss.NewStyle().Foreground(t.Subtext).Padding(0, 1).Render(stats)
	content.WriteString(statsStyled + "\n\n")

	// Ready beads list
	if len(m.ready) > 0 {
		content.WriteString(lipgloss.NewStyle().Foreground(t.Green).Bold(true).Padding(0, 1).Render("Top Ready") + "\n")
		for i, b := range m.ready {
			if i >= 5 {
				break
			}
			title := layout.TruncateRunes(b.Title, m.width-12, "...")
			line := fmt.Sprintf("  %s %s", b.ID, title)
			content.WriteString(lipgloss.NewStyle().Foreground(t.Text).Render(line) + "\n")
		}
	} else {
		content.WriteString("  No ready beads\n")
	}

	return content.String()
}
