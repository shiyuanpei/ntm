package palette

import (
	"fmt"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/tmux"
	"github.com/Dicklesworthstone/ntm/internal/tui/icons"
	"github.com/Dicklesworthstone/ntm/internal/tui/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionSelector is a TUI for selecting a tmux session
type SessionSelector struct {
	sessions []tmux.Session
	cursor   int
	selected string
	quitting bool
	width    int
	height   int

	// Theme
	theme theme.Theme
	icons icons.IconSet
}

// SessionSelectorKeyMap defines keybindings for the selector
type SessionSelectorKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
	Num1   key.Binding
	Num2   key.Binding
	Num3   key.Binding
	Num4   key.Binding
	Num5   key.Binding
	Num6   key.Binding
	Num7   key.Binding
	Num8   key.Binding
	Num9   key.Binding
}

var selectorKeys = SessionSelectorKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	Num1: key.NewBinding(key.WithKeys("1")),
	Num2: key.NewBinding(key.WithKeys("2")),
	Num3: key.NewBinding(key.WithKeys("3")),
	Num4: key.NewBinding(key.WithKeys("4")),
	Num5: key.NewBinding(key.WithKeys("5")),
	Num6: key.NewBinding(key.WithKeys("6")),
	Num7: key.NewBinding(key.WithKeys("7")),
	Num8: key.NewBinding(key.WithKeys("8")),
	Num9: key.NewBinding(key.WithKeys("9")),
}

// NewSessionSelector creates a new session selector
func NewSessionSelector(sessions []tmux.Session) SessionSelector {
	return SessionSelector{
		sessions: sessions,
		width:    60,
		height:   20,
		theme:    theme.Current(),
		icons:    icons.Current(),
	}
}

// Init implements tea.Model
func (s SessionSelector) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (s SessionSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, selectorKeys.Quit):
			s.quitting = true
			return s, tea.Quit

		case key.Matches(msg, selectorKeys.Up):
			if s.cursor > 0 {
				s.cursor--
			}

		case key.Matches(msg, selectorKeys.Down):
			if s.cursor < len(s.sessions)-1 {
				s.cursor++
			}

		case key.Matches(msg, selectorKeys.Select):
			if len(s.sessions) > 0 {
				s.selected = s.sessions[s.cursor].Name
				return s, tea.Quit
			}

		// Quick select with numbers 1-9
		case key.Matches(msg, selectorKeys.Num1):
			if s.selectByNumber(1) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num2):
			if s.selectByNumber(2) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num3):
			if s.selectByNumber(3) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num4):
			if s.selectByNumber(4) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num5):
			if s.selectByNumber(5) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num6):
			if s.selectByNumber(6) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num7):
			if s.selectByNumber(7) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num8):
			if s.selectByNumber(8) {
				return s, tea.Quit
			}
		case key.Matches(msg, selectorKeys.Num9):
			if s.selectByNumber(9) {
				return s, tea.Quit
			}
		}
	}

	return s, nil
}

func (s *SessionSelector) selectByNumber(n int) bool {
	idx := n - 1
	if idx >= 0 && idx < len(s.sessions) {
		s.cursor = idx
		s.selected = s.sessions[idx].Name
		return true
	}
	return false
}

// View implements tea.Model
func (s SessionSelector) View() string {
	t := s.theme
	ic := s.icons

	var b strings.Builder

	// Box width
	boxWidth := 50
	if s.width > 60 {
		boxWidth = 55
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary)

	b.WriteString("\n")
	b.WriteString(headerStyle.Render("  "+ic.Session+" Select Session") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render("  "+strings.Repeat("─", boxWidth)) + "\n\n")

	if len(s.sessions) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(t.Overlay).
			Italic(true)
		b.WriteString(emptyStyle.Render("  No tmux sessions found") + "\n")
		b.WriteString(emptyStyle.Render("  Create one with: ntm spawn <name>") + "\n\n")
	} else {
		for i, session := range s.sessions {
			isSelected := i == s.cursor

			var line strings.Builder

			// Cursor
			if isSelected {
				cursorStyle := lipgloss.NewStyle().Foreground(t.Pink).Bold(true)
				line.WriteString(cursorStyle.Render(ic.Pointer + " "))
			} else {
				line.WriteString("  ")
			}

			// Number (1-9)
			if i < 9 {
				numStyle := lipgloss.NewStyle().Foreground(t.Overlay)
				line.WriteString(numStyle.Render(fmt.Sprintf("%d ", i+1)))
			} else {
				line.WriteString("  ")
			}

			// Session name
			var nameStyle lipgloss.Style
			if isSelected {
				nameStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(t.Pink)
			} else {
				nameStyle = lipgloss.NewStyle().
					Foreground(t.Text)
			}
			line.WriteString(nameStyle.Render(session.Name))

			// Pane count and status
			infoStyle := lipgloss.NewStyle().Foreground(t.Subtext)
			info := fmt.Sprintf(" (%d windows)", session.Windows)
			line.WriteString(infoStyle.Render(info))

			// Attached indicator
			if session.Attached {
				attachedStyle := lipgloss.NewStyle().Foreground(t.Green)
				line.WriteString(attachedStyle.Render(" " + ic.Dot))
			}

			b.WriteString(line.String() + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render("  "+strings.Repeat("─", boxWidth)) + "\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(t.Overlay)
	keyStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)

	helpItems := []string{
		keyStyle.Render("↑/↓") + " navigate",
		keyStyle.Render("1-9") + " quick select",
		keyStyle.Render("enter") + " select",
		keyStyle.Render("esc") + " quit",
	}

	b.WriteString("  " + helpStyle.Render(strings.Join(helpItems, " • ")))

	return b.String()
}

// Selected returns the selected session name (empty if cancelled)
func (s SessionSelector) Selected() string {
	return s.selected
}

// RunSessionSelector runs the session selector and returns the selected session
func RunSessionSelector(sessions []tmux.Session) (string, error) {
	if len(sessions) == 0 {
		return "", fmt.Errorf("no sessions available")
	}

	// If only one session, return it directly
	if len(sessions) == 1 {
		return sessions[0].Name, nil
	}

	model := NewSessionSelector(sessions)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(SessionSelector)
	return result.Selected(), nil
}
