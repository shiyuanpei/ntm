package palette

import (
	"fmt"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/config"
	"github.com/Dicklesworthstone/ntm/internal/tmux"
	"github.com/Dicklesworthstone/ntm/internal/tui/icons"
	"github.com/Dicklesworthstone/ntm/internal/tui/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

// Phase represents the current UI phase
type Phase int

const (
	PhaseCommand Phase = iota
	PhaseTarget
	PhaseConfirm
)

// Target represents the send target
type Target int

const (
	TargetAll Target = iota
	TargetClaude
	TargetCodex
	TargetGemini
)

// Model is the Bubble Tea model for the palette
type Model struct {
	session   string
	commands  []config.PaletteCmd
	filtered  []config.PaletteCmd
	cursor    int
	selected  *config.PaletteCmd
	phase     Phase
	target    Target
	filter    textinput.Model
	width     int
	height    int
	sent      bool
	sentCount int
	quitting  bool
	err       error

	// visualOrder maps visual display position (0-indexed) to index in filtered slice.
	// This is needed because items are grouped by category, so visual order differs from slice order.
	visualOrder []int

	// Theme and styles
	theme  theme.Theme
	styles theme.Styles
	icons  icons.IconSet
}

// KeyMap defines the keybindings
type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Select  key.Binding
	Back    key.Binding
	Quit    key.Binding
	Filter  key.Binding
	Target1 key.Binding
	Target2 key.Binding
	Target3 key.Binding
	Target4 key.Binding
	Num1    key.Binding
	Num2    key.Binding
	Num3    key.Binding
	Num4    key.Binding
	Num5    key.Binding
	Num6    key.Binding
	Num7    key.Binding
	Num8    key.Binding
	Num9    key.Binding
}

var keys = KeyMap{
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
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back/quit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Target1: key.NewBinding(key.WithKeys("1")),
	Target2: key.NewBinding(key.WithKeys("2")),
	Target3: key.NewBinding(key.WithKeys("3")),
	Target4: key.NewBinding(key.WithKeys("4")),
	Num1:    key.NewBinding(key.WithKeys("1")),
	Num2:    key.NewBinding(key.WithKeys("2")),
	Num3:    key.NewBinding(key.WithKeys("3")),
	Num4:    key.NewBinding(key.WithKeys("4")),
	Num5:    key.NewBinding(key.WithKeys("5")),
	Num6:    key.NewBinding(key.WithKeys("6")),
	Num7:    key.NewBinding(key.WithKeys("7")),
	Num8:    key.NewBinding(key.WithKeys("8")),
	Num9:    key.NewBinding(key.WithKeys("9")),
}

// New creates a new palette model
func New(session string, commands []config.PaletteCmd) Model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter commands..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#f5c2e7"))

	t := theme.Current()
	s := theme.NewStyles(t)
	ic := icons.Current()

	m := Model{
		session:  session,
		commands: commands,
		filtered: commands,
		filter:   ti,
		phase:    PhaseCommand,
		width:    80,
		height:   24,
		theme:    t,
		styles:   s,
		icons:    ic,
	}

	// Build initial visual order mapping
	m.buildVisualOrder()

	return m
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.filter.Width = m.width/2 - 10
		return m, nil

	case tea.KeyMsg:
		switch m.phase {
		case PhaseCommand:
			return m.updateCommandPhase(msg)
		case PhaseTarget:
			return m.updateTargetPhase(msg)
		}
	}

	// Update filter input
	if m.phase == PhaseCommand {
		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.updateFiltered()
		return m, cmd
	}

	return m, nil
}

func (m *Model) updateCommandPhase(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		m.quitting = true
		return m, tea.Quit

	case key.Matches(msg, keys.Back):
		m.quitting = true
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if m.cursor > 0 {
			m.cursor--
		}

	case key.Matches(msg, keys.Down):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}

	case key.Matches(msg, keys.Select):
		if len(m.filtered) > 0 {
			m.selected = &m.filtered[m.cursor]
			m.phase = PhaseTarget
		}

	// Quick select with numbers 1-9
	case key.Matches(msg, keys.Num1):
		if m.selectByNumber(1) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num2):
		if m.selectByNumber(2) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num3):
		if m.selectByNumber(3) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num4):
		if m.selectByNumber(4) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num5):
		if m.selectByNumber(5) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num6):
		if m.selectByNumber(6) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num7):
		if m.selectByNumber(7) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num8):
		if m.selectByNumber(8) {
			m.phase = PhaseTarget
		}
	case key.Matches(msg, keys.Num9):
		if m.selectByNumber(9) {
			m.phase = PhaseTarget
		}

	default:
		// Let the textinput handle it
		var cmd tea.Cmd
		m.filter, cmd = m.filter.Update(msg)
		m.updateFiltered()
		return m, cmd
	}

	return m, nil
}

func (m *Model) selectByNumber(n int) bool {
	visualPos := n - 1 // Convert 1-based to 0-based
	if visualPos >= 0 && visualPos < len(m.visualOrder) {
		// Map visual position to actual index in filtered slice
		idx := m.visualOrder[visualPos]
		m.cursor = idx
		m.selected = &m.filtered[idx]
		return true
	}
	return false
}

func (m *Model) updateTargetPhase(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Back):
		m.phase = PhaseCommand
		m.selected = nil

	case key.Matches(msg, keys.Quit):
		m.quitting = true
		return m, tea.Quit

	case key.Matches(msg, keys.Target1):
		m.target = TargetAll
		return m.send()

	case key.Matches(msg, keys.Target2):
		m.target = TargetClaude
		return m.send()

	case key.Matches(msg, keys.Target3):
		m.target = TargetCodex
		return m.send()

	case key.Matches(msg, keys.Target4):
		m.target = TargetGemini
		return m.send()
	}

	return m, nil
}

func (m *Model) updateFiltered() {
	query := strings.ToLower(m.filter.Value())
	if query == "" {
		m.filtered = m.commands
	} else {
		m.filtered = nil
		for _, cmd := range m.commands {
			if strings.Contains(strings.ToLower(cmd.Label), query) ||
				strings.Contains(strings.ToLower(cmd.Key), query) ||
				strings.Contains(strings.ToLower(cmd.Category), query) {
				m.filtered = append(m.filtered, cmd)
			}
		}
	}

	// Build visual order mapping (items grouped by category)
	m.buildVisualOrder()

	// Keep cursor in bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// buildVisualOrder creates a mapping from visual position to filtered index.
// Items are grouped by category, so the visual order differs from the slice order.
func (m *Model) buildVisualOrder() {
	m.visualOrder = nil
	if len(m.filtered) == 0 {
		return
	}

	// Group by category (same logic as renderCommandList)
	categories := make(map[string][]int)
	categoryOrder := []string{}

	for i, cmd := range m.filtered {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		if _, exists := categories[cat]; !exists {
			categoryOrder = append(categoryOrder, cat)
		}
		categories[cat] = append(categories[cat], i)
	}

	// Build visual order following category grouping
	for _, cat := range categoryOrder {
		indices := categories[cat]
		m.visualOrder = append(m.visualOrder, indices...)
	}
}

func (m *Model) send() (tea.Model, tea.Cmd) {
	if m.selected == nil {
		return m, nil
	}

	panes, err := tmux.GetPanes(m.session)
	if err != nil {
		m.err = err
		return m, tea.Quit
	}

	prompt := m.selected.Prompt
	count := 0

	for _, p := range panes {
		var shouldSend bool

		switch m.target {
		case TargetAll:
			// Send to all agent panes
			shouldSend = p.Type != tmux.AgentUser
		case TargetClaude:
			shouldSend = p.Type == tmux.AgentClaude
		case TargetCodex:
			shouldSend = p.Type == tmux.AgentCodex
		case TargetGemini:
			shouldSend = p.Type == tmux.AgentGemini
		}

		if shouldSend {
			if err := tmux.SendKeys(p.ID, prompt, true); err != nil {
				m.err = err
				return m, tea.Quit
			}
			count++
		}
	}

	m.sent = true
	m.sentCount = count
	m.quitting = true
	return m, tea.Quit
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return m.viewQuitting()
	}

	switch m.phase {
	case PhaseCommand:
		return m.viewCommandPhase()
	case PhaseTarget:
		return m.viewTargetPhase()
	}

	return ""
}

func (m Model) viewQuitting() string {
	t := m.theme
	ic := m.icons

	if m.err != nil {
		errorStyle := lipgloss.NewStyle().Foreground(t.Error)
		return errorStyle.Render(fmt.Sprintf("%s Error: %v\n", ic.Cross, m.err))
	}

	if m.sent {
		successStyle := lipgloss.NewStyle().Foreground(t.Success)
		labelStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Text)

		targetName := "all agents"
		targetIcon := ic.All
		switch m.target {
		case TargetClaude:
			targetName = "Claude agents"
			targetIcon = ic.Claude
		case TargetCodex:
			targetName = "Codex agents"
			targetIcon = ic.Codex
		case TargetGemini:
			targetName = "Gemini agents"
			targetIcon = ic.Gemini
		}

		return successStyle.Render(ic.Check) + " " +
			labelStyle.Render(fmt.Sprintf("Sent \"%s\" to %s %s (%d panes)\n",
				m.selected.Label, targetIcon, targetName, m.sentCount))
	}

	return ""
}

func (m Model) viewCommandPhase() string {
	t := m.theme
	ic := m.icons

	var b strings.Builder

	// Calculate layout dimensions
	listWidth := m.width/2 - 2
	previewWidth := m.width/2 - 2
	if listWidth < 30 {
		listWidth = 30
	}
	if previewWidth < 30 {
		previewWidth = 30
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary).
		Padding(0, 1)

	sessionStyle := lipgloss.NewStyle().
		Foreground(t.Subtext)

	sepStyle := lipgloss.NewStyle().
		Foreground(t.Surface2)

	header := headerStyle.Render(ic.Palette+" Command Palette") +
		sepStyle.Render(" │ ") +
		sessionStyle.Render(m.session)

	b.WriteString(header + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render(strings.Repeat("─", m.width-4)) + "\n\n")

	// Filter input
	filterStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Surface1).
		Padding(0, 1).
		Width(listWidth - 4)

	b.WriteString("  " + filterStyle.Render(m.filter.View()) + "\n\n")

	// Build the two-column layout
	listContent := m.renderCommandList(listWidth - 4)
	previewContent := m.renderPreview(previewWidth - 4)

	// Create list box
	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Surface2).
		Width(listWidth - 2).
		Height(m.height - 12).
		Padding(1, 1)

	// Create preview box
	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Blue).
		Width(previewWidth - 2).
		Height(m.height - 12).
		Padding(1, 1)

	// Join horizontally
	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listBox.Render(listContent),
		"  ",
		previewBox.Render(previewContent),
	)

	b.WriteString(columns + "\n\n")

	// Help bar
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

func (m Model) renderCommandList(width int) string {
	t := m.theme
	ic := m.icons

	if len(m.filtered) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(t.Overlay).
			Italic(true)
		return emptyStyle.Render("No commands match your filter")
	}

	var lines []string

	// Group by category
	categories := make(map[string][]int)
	categoryOrder := []string{}

	for i, cmd := range m.filtered {
		cat := cmd.Category
		if cat == "" {
			cat = "General"
		}
		if _, exists := categories[cat]; !exists {
			categoryOrder = append(categoryOrder, cat)
		}
		categories[cat] = append(categories[cat], i)
	}

	itemNum := 0
	for _, cat := range categoryOrder {
		indices := categories[cat]

		// Category header
		catIcon := ic.CategoryIcon(cat)
		catStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Lavender)
		lines = append(lines, catStyle.Render(catIcon+" "+cat))

		for _, idx := range indices {
			cmd := m.filtered[idx]
			isSelected := idx == m.cursor
			itemNum++

			var line strings.Builder

			// Cursor
			if isSelected {
				cursorStyle := lipgloss.NewStyle().Foreground(t.Pink).Bold(true)
				line.WriteString(cursorStyle.Render(ic.Pointer + " "))
			} else {
				line.WriteString("  ")
			}

			// Number (1-9)
			if itemNum <= 9 {
				numStyle := lipgloss.NewStyle().Foreground(t.Overlay)
				line.WriteString(numStyle.Render(fmt.Sprintf("%d ", itemNum)))
			} else {
				line.WriteString("  ")
			}

			// Item
			var itemStyle lipgloss.Style
			if isSelected {
				itemStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(t.Pink)
			} else {
				itemStyle = lipgloss.NewStyle().
					Foreground(t.Text)
			}
			line.WriteString(itemStyle.Render(cmd.Label))

			lines = append(lines, line.String())
		}

		lines = append(lines, "") // Spacing between categories
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderPreview(width int) string {
	t := m.theme
	ic := m.icons

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary)

	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		emptyStyle := lipgloss.NewStyle().
			Foreground(t.Overlay).
			Italic(true)
		return emptyStyle.Render("Select a command to preview")
	}

	cmd := m.filtered[m.cursor]

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render(ic.Send+" "+cmd.Label) + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render(strings.Repeat("─", width)) + "\n\n")

	// Category badge
	if cmd.Category != "" {
		badgeStyle := lipgloss.NewStyle().
			Foreground(t.Base).
			Background(t.Mauve).
			Padding(0, 1)
		b.WriteString(badgeStyle.Render(cmd.Category) + "\n\n")
	}

	// Prompt content
	promptStyle := lipgloss.NewStyle().Foreground(t.Text)
	wrapped := wordwrap.String(cmd.Prompt, width)
	b.WriteString(promptStyle.Render(wrapped))

	return b.String()
}

func (m Model) viewTargetPhase() string {
	t := m.theme
	ic := m.icons

	var b strings.Builder

	// Box dimensions
	boxWidth := 60
	if m.width < 70 {
		boxWidth = m.width - 10
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Primary)

	b.WriteString("\n")
	b.WriteString(headerStyle.Render("  "+ic.Target+" Select Target") + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render("  "+strings.Repeat("─", boxWidth)) + "\n\n")

	// Selected command
	dimStyle := lipgloss.NewStyle().Foreground(t.Subtext)
	labelStyle := lipgloss.NewStyle().Foreground(t.Text)
	b.WriteString(dimStyle.Render("  Sending: ") + labelStyle.Render(m.selected.Label) + "\n\n")

	// Target options
	targets := []struct {
		key   string
		icon  string
		label string
		desc  string
		color lipgloss.Color
	}{
		{"1", ic.All, "All Agents", "broadcast to all agent panes", t.Green},
		{"2", ic.Claude, "Claude (cc)", "Anthropic agents only", t.Claude},
		{"3", ic.Codex, "Codex (cod)", "OpenAI agents only", t.Codex},
		{"4", ic.Gemini, "Gemini (gmi)", "Google agents only", t.Gemini},
	}

	for _, target := range targets {
		keyStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			Width(3)

		iconStyle := lipgloss.NewStyle().
			Foreground(target.color).
			Width(3)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Text).
			Width(16)

		descStyle := lipgloss.NewStyle().
			Foreground(t.Overlay)

		line := fmt.Sprintf("  %s %s %s %s\n",
			keyStyle.Render(target.key),
			iconStyle.Render(target.icon),
			labelStyle.Render(target.label),
			descStyle.Render(target.desc))

		b.WriteString(line)
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(t.Surface2).Render("  "+strings.Repeat("─", boxWidth)) + "\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().Foreground(t.Overlay)
	keyStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)

	helpItems := []string{
		keyStyle.Render("1-4") + " select target",
		keyStyle.Render("esc") + " back",
		keyStyle.Render("q") + " quit",
	}

	b.WriteString("  " + helpStyle.Render(strings.Join(helpItems, " • ")))

	return b.String()
}

// Result returns the send result after the program exits
func (m Model) Result() (sent bool, err error) {
	return m.sent, m.err
}
