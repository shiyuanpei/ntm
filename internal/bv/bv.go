// Package bv provides integration with the beads_viewer (bv) tool.
// It executes bv robot mode commands and parses their JSON output.
package bv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrNotInstalled indicates bv is not available
var ErrNotInstalled = errors.New("bv is not installed")

// ErrNoBaseline indicates no baseline exists for drift checking
var ErrNoBaseline = errors.New("no baseline found")

// IsInstalled checks if bv is available in PATH
func IsInstalled() bool {
	_, err := exec.LookPath("bv")
	return err == nil
}

// run executes bv with given args and returns stdout
func run(args ...string) (string, error) {
	if !IsInstalled() {
		return "", ErrNotInstalled
	}

	cmd := exec.Command("bv", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check for specific error conditions
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "No baseline found") {
			return "", ErrNoBaseline
		}
		return "", fmt.Errorf("bv %s: %w: %s", strings.Join(args, " "), err, stderrStr)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// GetInsights returns graph analysis insights (bottlenecks, keystones, etc.)
func GetInsights() (*InsightsResponse, error) {
	output, err := run("-robot-insights")
	if err != nil {
		return nil, err
	}

	var resp InsightsResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, fmt.Errorf("parsing insights: %w", err)
	}

	return &resp, nil
}

// GetPriority returns priority recommendations
func GetPriority() (*PriorityResponse, error) {
	output, err := run("-robot-priority")
	if err != nil {
		return nil, err
	}

	var resp PriorityResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, fmt.Errorf("parsing priority: %w", err)
	}

	return &resp, nil
}

// GetPlan returns a parallel execution plan
func GetPlan() (*PlanResponse, error) {
	output, err := run("-robot-plan")
	if err != nil {
		return nil, err
	}

	var resp PlanResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, fmt.Errorf("parsing plan: %w", err)
	}

	return &resp, nil
}

// GetRecipes returns available recipes
func GetRecipes() (*RecipesResponse, error) {
	output, err := run("-robot-recipes")
	if err != nil {
		return nil, err
	}

	var resp RecipesResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		return nil, fmt.Errorf("parsing recipes: %w", err)
	}

	return &resp, nil
}

// CheckDrift checks project drift from baseline
// Returns DriftResult with status and message
func CheckDrift() DriftResult {
	if !IsInstalled() {
		return DriftResult{
			Status:  DriftNoBaseline,
			Message: "bv not installed",
		}
	}

	cmd := exec.Command("bv", "-check-drift")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Parse exit code
	if err == nil {
		return DriftResult{
			Status:  DriftOK,
			Message: strings.TrimSpace(stdout.String()),
		}
	}

	// Check for exit code
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.ExitCode()
		message := strings.TrimSpace(stdout.String())
		if message == "" {
			message = strings.TrimSpace(stderr.String())
		}

		switch code {
		case 1:
			// Could be critical drift or no baseline
			if strings.Contains(message, "No baseline") {
				return DriftResult{
					Status:  DriftNoBaseline,
					Message: message,
				}
			}
			return DriftResult{
				Status:  DriftCritical,
				Message: message,
			}
		case 2:
			return DriftResult{
				Status:  DriftWarning,
				Message: message,
			}
		default:
			return DriftResult{
				Status:  DriftStatus(code),
				Message: message,
			}
		}
	}

	return DriftResult{
		Status:  DriftNoBaseline,
		Message: err.Error(),
	}
}

// GetTopBottlenecks returns the top N bottleneck issues
func GetTopBottlenecks(n int) ([]NodeScore, error) {
	insights, err := GetInsights()
	if err != nil {
		return nil, err
	}

	bottlenecks := insights.Bottlenecks
	if len(bottlenecks) > n {
		bottlenecks = bottlenecks[:n]
	}

	return bottlenecks, nil
}

// GetNextActions returns recommended next actions based on priority analysis
func GetNextActions(n int) ([]PriorityRecommendation, error) {
	priority, err := GetPriority()
	if err != nil {
		return nil, err
	}

	recommendations := priority.Recommendations
	if len(recommendations) > n {
		recommendations = recommendations[:n]
	}

	return recommendations, nil
}

// GetParallelTracks returns available parallel work tracks
func GetParallelTracks() ([]Track, error) {
	plan, err := GetPlan()
	if err != nil {
		return nil, err
	}

	return plan.Plan.Tracks, nil
}

// IsBottleneck checks if an issue ID is in the bottleneck list
func IsBottleneck(issueID string) (bool, float64, error) {
	insights, err := GetInsights()
	if err != nil {
		return false, 0, err
	}

	for _, b := range insights.Bottlenecks {
		if b.ID == issueID {
			return true, b.Value, nil
		}
	}

	return false, 0, nil
}

// HealthSummary returns a brief project health summary
type HealthSummary struct {
	DriftStatus   DriftStatus
	DriftMessage  string
	TopBottleneck string
	BottleneckCount int
}

// GetHealthSummary returns a quick project health check
func GetHealthSummary() (*HealthSummary, error) {
	summary := &HealthSummary{}

	// Check drift
	drift := CheckDrift()
	summary.DriftStatus = drift.Status
	summary.DriftMessage = drift.Message

	// Get bottlenecks
	bottlenecks, err := GetTopBottlenecks(5)
	if err != nil {
		// Non-fatal, just skip bottleneck info
		return summary, nil
	}

	summary.BottleneckCount = len(bottlenecks)
	if len(bottlenecks) > 0 {
		summary.TopBottleneck = bottlenecks[0].ID
	}

	return summary, nil
}
