package robot

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/agentmail"
	"github.com/Dicklesworthstone/ntm/internal/bv"
	"github.com/Dicklesworthstone/ntm/internal/tmux"
)

// DashboardOutput represents the comprehensive state for --robot-dashboard
type DashboardOutput struct {
	GeneratedAt time.Time `json:"generated_at"`
	Project     string    `json:"project"`
	Session     string    `json:"session,omitempty"`

	// Subsystem status
	Health HealthStatus   `json:"health"`
	Agents []AgentSummary `json:"agents"`
	Mail   MailSummary    `json:"mail"`
	Work   WorkSummary    `json:"work"` // beads
	Alerts []AlertSummary `json:"alerts"`
}

type HealthStatus struct {
	Status  string `json:"status"` // ok, warning, critical
	Message string `json:"message"`
}

type AgentSummary struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"` // idle, working, error
}

type MailSummary struct {
	Available   bool `json:"available"`
	UnreadCount int  `json:"unread_count"`
	UrgentCount int  `json:"urgent_count"`
}

type WorkSummary struct {
	ActiveBeads []string `json:"active_beads"`
	ReadyCount  int      `json:"ready_count"`
}

type AlertSummary struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// PrintDashboard gathers and outputs a high-level project dashboard.
func PrintDashboard() error {
	cwd, _ := os.Getwd()
	output := DashboardOutput{
		GeneratedAt: time.Now().UTC(),
		Project:     filepath.Base(cwd),
	}

	// 1. Session & Agents
	if tmux.InTmux() {
		output.Session = tmux.GetCurrentSession()
		panes, err := tmux.GetPanes(output.Session)
		if err == nil {
			for _, p := range panes {
				if p.Type == tmux.AgentUser {
					continue
				}
				output.Agents = append(output.Agents, AgentSummary{
					Name:   p.Title,
					Type:   string(p.Type),
					Status: "unknown", 
				})
			}
		}
	}

	// 2. Health (bv)
	if bv.IsInstalled() {
		drift := bv.CheckDrift()
		output.Health = HealthStatus{
			Status:  string(drift.Status),
			Message: drift.Message,
		}
	} else {
		output.Health = HealthStatus{Status: "unavailable", Message: "bv not installed"}
	}

	// 3. Mail
	client := agentmail.NewClient(agentmail.WithProjectKey(cwd))
	if client.IsAvailable() {
		output.Mail.Available = true
	}

	return encodeJSON(output)
}
