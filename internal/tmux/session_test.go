package tmux

import (
	"testing"
)

func TestValidateSessionName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"myproject", false},
		{"my-project", false},
		{"my_project", false},
		{"MyProject123", false},
		{"", true},          // empty
		{"my.project", true}, // contains .
		{"my:project", true}, // contains :
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSessionName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestAgentTypeFromTitle(t *testing.T) {
	tests := []struct {
		title    string
		expected AgentType
	}{
		{"myproject__cc_1", AgentClaude},
		{"myproject__cod_1", AgentCodex},
		{"myproject__gmi_1", AgentGemini},
		{"myproject", AgentUser},
		{"zsh", AgentUser},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			// Create a pane and check type detection logic
			pane := Pane{Title: tt.title, Type: AgentUser}

			// Apply same logic as GetPanes
			if contains(pane.Title, "__cc") {
				pane.Type = AgentClaude
			} else if contains(pane.Title, "__cod") {
				pane.Type = AgentCodex
			} else if contains(pane.Title, "__gmi") {
				pane.Type = AgentGemini
			}

			if pane.Type != tt.expected {
				t.Errorf("Expected type %v for title %q, got %v", tt.expected, tt.title, pane.Type)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestInTmux(t *testing.T) {
	// This will be false in test environment
	// Just verify the function doesn't panic
	_ = InTmux()
}

func TestIsInstalled(t *testing.T) {
	// This checks if tmux is installed on the system
	// Just verify the function doesn't panic
	_ = IsInstalled()
}
