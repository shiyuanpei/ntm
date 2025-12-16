package components

import (
	"strings"
	"testing"
)

func TestRenderEmptyState(t *testing.T) {
	t.Run("basic rendering with all fields", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:        IconWaiting,
			Title:       "No metrics yet",
			Description: "Data appears when agents start",
			Action:      "Press r to refresh",
			Width:       40,
			Centered:    true,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
		if !strings.Contains(out, "No metrics yet") {
			t.Errorf("expected title in output, got %q", out)
		}
		if !strings.Contains(out, "Data appears") {
			t.Errorf("expected description in output, got %q", out)
		}
		if !strings.Contains(out, "Press r") {
			t.Errorf("expected action in output, got %q", out)
		}
	})

	t.Run("minimal rendering (title only)", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:     IconEmpty,
			Title:    "Nothing found",
			Width:    30,
			Centered: true,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
		if !strings.Contains(out, "Nothing found") {
			t.Errorf("expected title in output, got %q", out)
		}
	})

	t.Run("default title when empty", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:     IconEmpty,
			Width:    30,
			Centered: true,
		})
		if !strings.Contains(out, "Nothing to show") {
			t.Errorf("expected default title, got %q", out)
		}
	})

	t.Run("success icon styling", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:        IconSuccess,
			Title:       "All clear",
			Description: "No alerts",
			Width:       30,
			Centered:    true,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
		if !strings.Contains(out, "All clear") {
			t.Errorf("expected title in output, got %q", out)
		}
	})

	t.Run("external icon for external action needed", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:        IconExternal,
			Title:       "Not initialized",
			Description: "Run 'bd init' in your project",
			Width:       40,
			Centered:    true,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
		if !strings.Contains(out, "Not initialized") {
			t.Errorf("expected title in output, got %q", out)
		}
	})

	t.Run("unknown icon fallback", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:     IconUnknown,
			Title:    "Unknown state",
			Width:    30,
			Centered: true,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
	})

	t.Run("truncates long description", func(t *testing.T) {
		longDesc := "This is a very very very long description that should be truncated"
		out := RenderEmptyState(EmptyStateOptions{
			Icon:        IconWaiting,
			Title:       "Test",
			Description: longDesc,
			Width:       20,
			Centered:    true,
		})
		if !strings.Contains(out, "…") {
			t.Errorf("expected truncation ellipsis, got %q", out)
		}
	})

	t.Run("left-aligned mode", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:     IconEmpty,
			Title:    "No items",
			Width:    30,
			Centered: false,
		})
		if out == "" {
			t.Fatal("expected non-empty output")
		}
		// Left-aligned should have padding
		if !strings.Contains(out, "No items") {
			t.Errorf("expected title in output, got %q", out)
		}
	})

	t.Run("zero width renders without crash", func(t *testing.T) {
		out := RenderEmptyState(EmptyStateOptions{
			Icon:  IconEmpty,
			Title: "Test",
			Width: 0,
		})
		if out == "" {
			t.Fatal("expected non-empty output even with zero width")
		}
	})
}

func TestEmptyStateIcons(t *testing.T) {
	icons := []EmptyStateIcon{
		IconWaiting,
		IconEmpty,
		IconExternal,
		IconSuccess,
		IconUnknown,
	}

	for _, icon := range icons {
		t.Run(string(icon), func(t *testing.T) {
			resolved := resolveEmptyIcon(icon)
			if resolved == "" {
				t.Errorf("icon %q resolved to empty string", icon)
			}
		})
	}

	t.Run("invalid icon uses fallback", func(t *testing.T) {
		resolved := resolveEmptyIcon(EmptyStateIcon("invalid"))
		if resolved == "" {
			t.Error("invalid icon should fallback to info icon, got empty string")
		}
	})
}
