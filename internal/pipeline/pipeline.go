package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/status"
	"github.com/Dicklesworthstone/ntm/internal/tmux"
)

// Stage represents a step in the pipeline
type Stage struct {
	AgentType string
	Prompt    string
	Model     string // Optional
}

// Pipeline represents a sequence of stages
type Pipeline struct {
	Session string
	Stages  []Stage
}

// Execute runs the pipeline stages sequentially
func Execute(ctx context.Context, p Pipeline) error {
	var previousOutput string
	var lastPaneID string

	detector := status.NewDetector()

	for i, stage := range p.Stages {
		fmt.Printf("Stage %d/%d [%s]: %s\n", i+1, len(p.Stages), stage.AgentType, truncate(stage.Prompt, 50))

		// 1. Find a suitable pane
		paneID, err := findPaneForStage(p.Session, stage.AgentType, stage.Model)
		if err != nil {
			return fmt.Errorf("stage %d failed: %w", i+1, err)
		}

		// Capture state BEFORE sending prompt to isolate the new response later
		beforeOutput, err := tmux.CapturePaneOutput(paneID, 2000)
		if err != nil {
			// Non-fatal, just means we might capture too much context
			beforeOutput = ""
		}

		// 2. Prepare prompt
		prompt := stage.Prompt
		if previousOutput != "" {
			if paneID == lastPaneID {
				// Same agent, context is already in the pane history.
				// Add a reference to the previous output instead of duplicating it.
				prompt = fmt.Sprintf("%s\n\n(See previous output above)", prompt)
			} else {
				// Different agent, inject the output from the previous stage
				prompt = fmt.Sprintf("%s\n\nResult from previous stage:\n%s", prompt, previousOutput)
			}
		}

		// 3. Send prompt
		if err := tmux.SendKeys(paneID, prompt, true); err != nil {
			return fmt.Errorf("stage %d sending prompt: %w", i+1, err)
		}

		// 4. Wait for working state (debounce)
		time.Sleep(2 * time.Second)

		// 5. Wait for idle state
		fmt.Printf("  Waiting for agent...")
		if err := waitForIdle(ctx, detector, paneID); err != nil {
			return fmt.Errorf("stage %d waiting for completion: %w", i+1, err)
		}
		fmt.Println(" Done.")

		// 6. Capture output
		// We capture a larger buffer to ensure we get the full response.
		// 2000 lines should cover most responses without being excessive.
		afterOutput, err := tmux.CapturePaneOutput(paneID, 2000)
		if err != nil {
			return fmt.Errorf("stage %d capturing output: %w", i+1, err)
		}

		// Extract only the new content
		output := extractNewOutput(beforeOutput, afterOutput)
		previousOutput = output
		lastPaneID = paneID
	}

	return nil
}

// extractNewOutput isolates the text added to the pane after the before state.
func extractNewOutput(before, after string) string {
	if before == "" {
		return after
	}
	if after == "" {
		return ""
	}

	// Simple suffix check: if 'after' ends with 'before' (unlikely as new text is appended),
	// or 'after' starts with 'before' (likely).
	if len(after) >= len(before) && after[:len(before)] == before {
		return after[len(before):]
	}

	// If history scrolled (before is not a strict prefix), try to overlap.
	// Find the longest suffix of 'before' that matches a prefix of 'after'.
	// This is O(N^2) in worst case but N=2000 chars/lines is small.
	// Optimization: limit search window.
	
	// Actually, tmux capture is usually stable. If buffer shifted, 'before' start is gone.
	// But 'after' should contain the tail of 'before'.
	// Find where 'before' ends inside 'after'.
	
	// Heuristic: take the last 100 chars of 'before' and find them in 'after'.
	const overlapSize = 100
	if len(before) > overlapSize {
		tail := before[len(before)-overlapSize:]
		if idx := strings.Index(after, tail); idx != -1 {
			// Found the overlap point
			return after[idx+len(tail):]
		}
	} else {
		if idx := strings.Index(after, before); idx != -1 {
			return after[idx+len(before):]
		}
	}

	// Fallback: return everything if we can't match (better than nothing)
	return after
}

func findPaneForStage(session, agentType, model string) (string, error) {
	panes, err := tmux.GetPanes(session)
	if err != nil {
		return "", err
	}

	// First pass: look for exact match (type + model)
	for _, p := range panes {
		if string(p.Type) == agentType {
			// Check model if specified
			if model != "" && p.Variant != model {
				continue
			}
			return p.ID, nil
		}
	}

	// Second pass: relaxed match (type only, ignore model mismatch if not found)
	// Only if model was specified but not found
	if model != "" {
		for _, p := range panes {
			if string(p.Type) == agentType {
				return p.ID, nil
			}
		}
	}

	return "", fmt.Errorf("no agent found for type %s (model %s)", agentType, model)
}

func waitForIdle(ctx context.Context, detector status.Detector, paneID string) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.After(30 * time.Minute) // Max 30 min per stage default

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for agent")
		case <-ticker.C:
			s, err := detector.Detect(paneID)
			if err != nil {
				continue
			}
			if s.State == status.StateIdle {
				return nil
			}
			// Optional: print progress indicator
		}
	}
}

func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
