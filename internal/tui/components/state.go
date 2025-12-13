package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Dicklesworthstone/ntm/internal/tui/icons"
	"github.com/Dicklesworthstone/ntm/internal/tui/layout"
	"github.com/Dicklesworthstone/ntm/internal/tui/theme"
)

type StateKind int

const (
	StateEmpty StateKind = iota
	StateLoading
	StateError
)

type StateOptions struct {
	Kind    StateKind
	Icon    string
	Message string
	Hint    string
	Width   int
	Align   lipgloss.Position
}

func RenderState(opts StateOptions) string {
	t := theme.Current()
	ic := icons.Current()

	align := opts.Align
	indent := "  "
	if align == lipgloss.Center {
		indent = ""
	}

	icon := strings.TrimSpace(opts.Icon)
	lineStyle := lipgloss.NewStyle().Foreground(t.Overlay).Italic(true)
	hintStyle := lipgloss.NewStyle().Foreground(t.Overlay).Italic(true)

	message := strings.TrimSpace(opts.Message)
	hint := strings.TrimSpace(opts.Hint)

	switch opts.Kind {
	case StateLoading:
		lineStyle = lipgloss.NewStyle().Foreground(t.Subtext).Italic(true)
		if message == "" {
			message = "Loading…"
		}
		if icon == "" {
			icon = strings.TrimSpace(ic.Gear)
			if icon == "" {
				icon = "…"
			}
		}
	case StateError:
		lineStyle = lipgloss.NewStyle().Foreground(t.Red).Italic(true)
		hintStyle = lipgloss.NewStyle().Foreground(t.Overlay).Italic(true)
		if message == "" {
			message = "Something went wrong"
		}
		if icon == "" {
			icon = strings.TrimSpace(ic.Warning)
			if icon == "" {
				icon = "!"
			}
		}
	default:
		if message == "" {
			message = "Nothing to show"
		}
		if icon == "" {
			icon = strings.TrimSpace(ic.Info)
			if icon == "" {
				icon = "i"
			}
		}
	}

	width := opts.Width
	if width < 0 {
		width = 0
	}

	prefix := indent + icon
	if icon != "" {
		prefix += " "
	}

	available := width
	if available > 0 {
		available -= lipgloss.Width(prefix)
		if available < 0 {
			available = 0
		}
	}

	if available > 0 {
		message = layout.TruncateRunes(message, available, "…")
	}

	lines := []string{lineStyle.Render(prefix + message)}

	if hint != "" {
		hintPrefix := indent
		hAvailable := width
		if hAvailable > 0 {
			hAvailable -= lipgloss.Width(hintPrefix)
			if hAvailable < 0 {
				hAvailable = 0
			}
		}
		if hAvailable > 0 {
			hint = layout.TruncateRunes(hint, hAvailable, "…")
		}
		lines = append(lines, hintStyle.Render(hintPrefix+hint))
	}

	rendered := strings.Join(lines, "\n")
	if width > 0 && (align == lipgloss.Center || align == lipgloss.Right) {
		return lipgloss.NewStyle().Width(width).Align(align).Render(rendered)
	}

	return rendered
}

func EmptyState(message string, width int) string {
	return RenderState(StateOptions{Kind: StateEmpty, Message: message, Width: width})
}

func LoadingState(message string, width int) string {
	return RenderState(StateOptions{Kind: StateLoading, Message: message, Width: width})
}

func ErrorState(message, hint string, width int) string {
	return RenderState(StateOptions{Kind: StateError, Message: message, Hint: hint, Width: width})
}
