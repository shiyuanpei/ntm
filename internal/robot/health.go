// Package robot provides machine-readable output for AI agents.
// health.go contains the --robot-health flag implementation.
package robot

import (
	"time"

	"github.com/Dicklesworthstone/ntm/internal/bv"
)

// HealthOutput provides a focused project health summary for AI agents
type HealthOutput struct {
	GeneratedAt       time.Time              `json:"generated_at"`
	BvAvailable       bool                   `json:"bv_available"`
	BdAvailable       bool                   `json:"bd_available"`
	Error             string                 `json:"error,omitempty"`
	DriftStatus       string                 `json:"drift_status,omitempty"`
	DriftMessage      string                 `json:"drift_message,omitempty"`
	TopBottlenecks    []bv.NodeScore         `json:"top_bottlenecks,omitempty"`
	TopKeystones      []bv.NodeScore         `json:"top_keystones,omitempty"`
	ReadyCount        int                    `json:"ready_count"`
	InProgressCount   int                    `json:"in_progress_count"`
	BlockedCount      int                    `json:"blocked_count"`
	NextRecommended   []RecommendedAction    `json:"next_recommended,omitempty"`
	DependencyContext *bv.DependencyContext  `json:"dependency_context,omitempty"`
}

// RecommendedAction is a simplified priority recommendation
type RecommendedAction struct {
	IssueID  string `json:"issue_id"`
	Title    string `json:"title"`
	Reason   string `json:"reason"`
	Priority int    `json:"priority"`
}

// PrintHealth outputs a focused project health summary for AI consumption
func PrintHealth() error {
	output := HealthOutput{
		GeneratedAt: time.Now().UTC(),
		BvAvailable: bv.IsInstalled(),
		BdAvailable: bv.IsBdInstalled(),
	}

	// Get drift status
	drift := bv.CheckDrift()
	output.DriftStatus = drift.Status.String()
	output.DriftMessage = drift.Message

	// Get top bottlenecks (limit to 5)
	bottlenecks, err := bv.GetTopBottlenecks(5)
	if err == nil {
		output.TopBottlenecks = bottlenecks
	}

	// Get insights for keystones
	insights, err := bv.GetInsights()
	if err == nil && insights != nil {
		keystones := insights.Keystones
		if len(keystones) > 5 {
			keystones = keystones[:5]
		}
		output.TopKeystones = keystones
	}

	// Get priority recommendations
	recommendations, err := bv.GetNextActions(5)
	if err == nil {
		for _, rec := range recommendations {
			var reason string
			if len(rec.Reasoning) > 0 {
				reason = rec.Reasoning[0]
			}
			output.NextRecommended = append(output.NextRecommended, RecommendedAction{
				IssueID:  rec.IssueID,
				Title:    rec.Title,
				Reason:   reason,
				Priority: rec.SuggestedPriority,
			})
		}
	}

	// Get dependency context (includes ready/in-progress/blocked counts)
	depCtx, err := bv.GetDependencyContext(5)
	if err == nil {
		output.DependencyContext = depCtx
		output.ReadyCount = depCtx.ReadyCount
		output.BlockedCount = depCtx.BlockedCount
		output.InProgressCount = len(depCtx.InProgressTasks)
	}

	return encodeJSON(output)
}
