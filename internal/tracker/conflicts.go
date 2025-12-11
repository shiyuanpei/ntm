package tracker

import (
	"time"
)

// Conflict represents a detected file conflict
type Conflict struct {
	Path     string               `json:"path"`
	Changes  []RecordedFileChange `json:"changes,omitempty"`
	Severity string               `json:"severity,omitempty"` // "warning", "critical"
	Agents   []string             `json:"agents,omitempty"`
	LastAt   time.Time            `json:"last_at,omitempty"`
}

// DetectConflicts analyzes a set of changes for conflicts.
func DetectConflicts(changes []RecordedFileChange) []Conflict {
	// Group by file path
	byPath := make(map[string][]RecordedFileChange)
	for _, change := range changes {
		// Only care about modifications for now
		if change.Change.Type == FileModified {
			byPath[change.Change.Path] = append(byPath[change.Change.Path], change)
		}
	}

	var conflicts []Conflict
	for path, pathChanges := range byPath {
		if len(pathChanges) > 1 {
			// Check if different agents involved
			// We consider it a conflict if the sets of agents differ or if distinct agents touched it
			// For simplicity: if we have more than one change event, and the agents involved are not identical
			// (or even if they are, multiple writes to same file might be race condition?)
			// Let's stick to "Modified by different agents" heuristic.

			// Collect all unique agents involved across all changes
			allAgents := make(map[string]bool)
			for _, pc := range pathChanges {
				for _, agent := range pc.Agents {
					allAgents[agent] = true
				}
			}

			// If only 1 agent ever touched it, no conflict (unless it's a self-overwrite race?)
			// But usually conflict implies >= 2 actors.
			if len(allAgents) > 1 {
				agentList := make([]string, 0, len(allAgents))
				for agent := range allAgents {
					agentList = append(agentList, agent)
				}
				
				conflicts = append(conflicts, Conflict{
					Path:     path,
					Changes:  pathChanges,
					Severity: "warning", // Default to warning for same-file
					Agents:   agentList,
				})
			}
		}
	}
	return conflicts
}

// DetectConflictsRecent analyzes global file changes within the given window.
func DetectConflictsRecent(window time.Duration) []Conflict {
	changes := GlobalFileChanges.Since(time.Now().Add(-window))
	return DetectConflicts(changes)
}

// ConflictsSince returns files changed by more than one agent since the timestamp.
func ConflictsSince(ts time.Time, session string) []Conflict {
	changes := GlobalFileChanges.Since(ts)
	if len(changes) == 0 {
		return nil
	}

	byPath := make(map[string]map[string]time.Time)
	for _, c := range changes {
		if session != "" && c.Session != session {
			continue
		}
		if len(c.Agents) == 0 {
			continue
		}
		if _, ok := byPath[c.Change.Path]; !ok {
			byPath[c.Change.Path] = make(map[string]time.Time)
		}
		for _, agent := range c.Agents {
			byPath[c.Change.Path][agent] = c.Timestamp
		}
	}

	out := make([]Conflict, 0)
	for path, agentsMap := range byPath {
		if len(agentsMap) < 2 {
			continue
		}
		conflict := Conflict{Path: path}
		var last time.Time
		for agent, ts := range agentsMap {
			conflict.Agents = append(conflict.Agents, agent)
			if ts.After(last) {
				last = ts
			}
		}
		conflict.LastAt = last
		conflict.Severity = "warning"
		out = append(out, conflict)
	}
	return out
}