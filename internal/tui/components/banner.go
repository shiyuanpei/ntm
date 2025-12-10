package components

import (
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// ASCII art logos for NTM
var (
	// Large banner logo
	LogoLarge = []string{
		"███╗   ██╗████████╗███╗   ███╗",
		"████╗  ██║╚══██╔══╝████╗ ████║",
		"██╔██╗ ██║   ██║   ██╔████╔██║",
		"██║╚██╗██║   ██║   ██║╚██╔╝██║",
		"██║ ╚████║   ██║   ██║ ╚═╝ ██║",
		"╚═╝  ╚═══╝   ╚═╝   ╚═╝     ╚═╝",
	}

	// Medium banner logo
	LogoMedium = []string{
		"╔╗╔╔╦╗╔╦╗",
		"║║║ ║ ║║║",
		"╝╚╝ ╩ ╩ ╩",
	}

	// Small inline logo
	LogoSmall = "⟦NTM⟧"

	// Icon variants
	LogoIcon      = "󰆍" // Terminal icon (Nerd Font)
	LogoIconPlain = "▣"  // Plain Unicode fallback
)

// Catppuccin gradient colors
var (
	GradientPrimary = []string{
		"#89b4fa", // Blue
		"#b4befe", // Lavender
		"#cba6f7", // Mauve
	}

	GradientSecondary = []string{
		"#cba6f7", // Mauve
		"#f5c2e7", // Pink
		"#f38ba8", // Red
	}

	GradientSuccess = []string{
		"#94e2d5", // Teal
		"#a6e3a1", // Green
		"#f9e2af", // Yellow
	}

	GradientRainbow = []string{
		"#f38ba8", // Red
		"#fab387", // Peach
		"#f9e2af", // Yellow
		"#a6e3a1", // Green
		"#89dceb", // Sky
		"#89b4fa", // Blue
		"#cba6f7", // Mauve
	}

	GradientAgent = map[string][]string{
		"claude": {"#cba6f7", "#b4befe", "#89b4fa"}, // Purple to blue
		"codex":  {"#89b4fa", "#74c7ec", "#89dceb"}, // Blue to cyan
		"gemini": {"#f9e2af", "#fab387", "#f38ba8"}, // Yellow to red
	}
)

// RenderBanner renders the large logo with gradient
func RenderBanner(animated bool, tick int) string {
	var lines []string

	for _, line := range LogoLarge {
		if animated {
			lines = append(lines, styles.Shimmer(line, tick, GradientPrimary...))
		} else {
			lines = append(lines, styles.GradientText(line, GradientPrimary...))
		}
	}

	return strings.Join(lines, "\n")
}

// RenderBannerMedium renders the medium logo with gradient
func RenderBannerMedium(animated bool, tick int) string {
	var lines []string

	for _, line := range LogoMedium {
		if animated {
			lines = append(lines, styles.Shimmer(line, tick, GradientPrimary...))
		} else {
			lines = append(lines, styles.GradientText(line, GradientPrimary...))
		}
	}

	return strings.Join(lines, "\n")
}

// RenderInlineLogo renders a small inline logo
func RenderInlineLogo() string {
	return styles.GradientText(LogoSmall, GradientPrimary...)
}

// RenderSubtitle renders a styled subtitle
func RenderSubtitle(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6adc8")).
		Italic(true).
		Render(text)
}

// RenderVersion renders a styled version string
func RenderVersion(version string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Render("v" + version)
}

// RenderHeaderBar renders a full header bar with title
func RenderHeaderBar(title string, width int) string {
	// Gradient divider
	divider := styles.GradientDivider(width, GradientPrimary...)

	// Centered title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#cdd6f4"))

	centeredTitle := styles.CenterText(titleStyle.Render(title), width)

	return divider + "\n" + centeredTitle + "\n" + divider
}

// RenderSection renders a section header
func RenderSection(title string, width int) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#89b4fa"))

	// Gradient line after title
	titleLen := len(title) + 2
	remaining := width - titleLen
	if remaining < 0 {
		remaining = 0
	}

	line := styles.GradientText(strings.Repeat("─", remaining), GradientPrimary...)

	return titleStyle.Render(title) + " " + line
}

// RenderAgentBadge renders a colored badge for an agent type
func RenderAgentBadge(agentType string) string {
	var bgColor, fgColor string
	var icon string

	switch agentType {
	case "claude", "cc":
		bgColor = "#cba6f7"
		fgColor = "#1e1e2e"
		icon = "󰗣"
	case "codex", "cod":
		bgColor = "#89b4fa"
		fgColor = "#1e1e2e"
		icon = ""
	case "gemini", "gmi":
		bgColor = "#f9e2af"
		fgColor = "#1e1e2e"
		icon = "󰊤"
	default:
		bgColor = "#a6e3a1"
		fgColor = "#1e1e2e"
		icon = ""
	}

	return lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor)).
		Bold(true).
		Padding(0, 1).
		Render(icon + " " + strings.ToUpper(agentType))
}

// RenderStatusBadge renders a status badge
func RenderStatusBadge(status string) string {
	var bgColor string
	var icon string

	switch status {
	case "running", "active":
		bgColor = "#a6e3a1"
		icon = "●"
	case "idle":
		bgColor = "#f9e2af"
		icon = "○"
	case "error", "failed":
		bgColor = "#f38ba8"
		icon = "✗"
	case "success", "done":
		bgColor = "#a6e3a1"
		icon = "✓"
	default:
		bgColor = "#6c7086"
		icon = "•"
	}

	return lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color("#1e1e2e")).
		Padding(0, 1).
		Render(icon + " " + status)
}

// RenderKeyMap renders a keyboard shortcuts help section
func RenderKeyMap(keys map[string]string, width int) string {
	var lines []string

	keyStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#45475a")).
		Foreground(lipgloss.Color("#cdd6f4")).
		Bold(true).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6adc8"))

	for key, desc := range keys {
		lines = append(lines, keyStyle.Render(key)+" "+descStyle.Render(desc))
	}

	// Join with separator
	return strings.Join(lines, "  ")
}

// RenderFooter renders a styled footer
func RenderFooter(text string, width int) string {
	divider := styles.GradientDivider(width, "#45475a", "#313244")

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true)

	return divider + "\n" + styles.CenterText(footerStyle.Render(text), width)
}

// RenderHint renders a dimmed hint text
func RenderHint(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Italic(true).
		Render(text)
}

// RenderHighlight renders highlighted text
func RenderHighlight(text string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5e0dc")).
		Bold(true).
		Render(text)
}

// RenderCommand renders a command with styling
func RenderCommand(cmd string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#89b4fa")).
		Bold(true).
		Render(cmd)
}

// RenderArg renders an argument with styling
func RenderArg(arg string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Render("<" + arg + ">")
}

// RenderFlag renders a flag with styling
func RenderFlag(flag string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f9e2af")).
		Render(flag)
}

// RenderExample renders an example with styling
func RenderExample(example string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fab387")).
		Italic(true).
		Render(example)
}
