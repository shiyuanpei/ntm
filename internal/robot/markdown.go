package robot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/bv"
)

// AgentTable renders a markdown table summarizing agents per session.
func AgentTable(sessions []SnapshotSession) string {
	var b strings.Builder
	b.WriteString("| Session | Pane | Type | Variant | State |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, sess := range sessions {
		for _, agent := range sess.Agents {
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n",
				sess.Name,
				agent.Pane,
				agent.Type,
				agent.Variant,
				agent.State)
		}
	}
	return b.String()
}

// AlertsList renders alerts as a markdown bullet list.
func AlertsList(alerts []AlertInfo) string {
	if len(alerts) == 0 {
		return "_No active alerts._"
	}
	var b strings.Builder
	for _, a := range alerts {
		fmt.Fprintf(&b, "- [%s] %s", strings.ToUpper(a.Severity), a.Message)
		if a.Session != "" {
			fmt.Fprintf(&b, " (session: %s", a.Session)
			if a.Pane != "" {
				fmt.Fprintf(&b, ", pane: %s", a.Pane)
			}
			fmt.Fprintf(&b, ")")
		}
		if a.BeadID != "" {
			fmt.Fprintf(&b, " [bead: %s]", a.BeadID)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// BeadsSummary renders a concise markdown summary of bead counts.
func BeadsSummary(summary *bv.BeadsSummary) string {
	if summary == nil || !summary.Available {
		return "_Beads summary unavailable._"
	}
	return fmt.Sprintf(
		"- Total: %d (Open: %d, In Progress: %d, Blocked: %d, Ready: %d, Closed: %d)",
		summary.Total,
		summary.Open,
		summary.InProgress,
		summary.Blocked,
		summary.Ready,
		summary.Closed,
	)
}

// SuggestedActions renders planned actions as markdown list items.
func SuggestedActions(actions []BeadAction) string {
	if len(actions) == 0 {
		return "_No suggested actions._"
	}
	var b strings.Builder
	for _, act := range actions {
		fmt.Fprintf(&b, "- %s: %s", act.BeadID, act.Title)
		if len(act.BlockedBy) > 0 {
			fmt.Fprintf(&b, " (blocked by: %s)", strings.Join(act.BlockedBy, ", "))
		}
		if act.Command != "" {
			fmt.Fprintf(&b, " — `%s`", act.Command)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// AgentTableRow represents a row in the agent markdown table.
type AgentTableRow struct {
	Agent  string
	Type   string
	Status string
}

// SuggestedAction is a lightweight action item for numbered lists.
type SuggestedAction struct {
	Title  string
	Reason string
}

// RenderAgentTable returns a markdown table of agents.
func RenderAgentTable(rows []AgentTableRow) string {
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("| Agent | Type | Status |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", r.Agent, r.Type, r.Status)
	}
	return b.String()
}

// RenderAlertsList groups alerts by severity and returns markdown bullets.
// Order of severities is: critical, warning, info, other.
func RenderAlertsList(alerts []AlertInfo) string {
	if len(alerts) == 0 {
		return ""
	}

	grouped := make(map[string][]AlertInfo)
	for _, a := range alerts {
		sev := strings.ToLower(a.Severity)
		grouped[sev] = append(grouped[sev], a)
	}

	severityOrder := []string{"critical", "warning", "info"}

	var b strings.Builder
	for _, sev := range severityOrder {
		if len(grouped[sev]) == 0 {
			continue
		}
		fmt.Fprintf(&b, "### %s\n", strings.Title(sev))
		for _, a := range grouped[sev] {
			loc := strings.TrimSpace(strings.Join([]string{a.Session, a.Pane}, " "))
			if loc != "" {
				loc = " (" + loc + ")"
			}
			fmt.Fprintf(&b, "- [%s] %s%s\n", a.Type, a.Message, loc)
		}
		b.WriteString("\n")
	}

	var others []string
	for sev := range grouped {
		if sev != "critical" && sev != "warning" && sev != "info" {
			others = append(others, sev)
		}
	}
	sort.Strings(others)
	for _, sev := range others {
		fmt.Fprintf(&b, "### %s\n", strings.Title(sev))
		for _, a := range grouped[sev] {
			fmt.Fprintf(&b, "- [%s] %s\n", a.Type, a.Message)
		}
		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

// RenderSuggestedActions returns a numbered markdown list.
func RenderSuggestedActions(actions []SuggestedAction) string {
	if len(actions) == 0 {
		return ""
	}
	var b strings.Builder
	for i, a := range actions {
		line := a.Title
		if a.Reason != "" {
			line = fmt.Sprintf("%s — %s", a.Title, a.Reason)
		}
		fmt.Fprintf(&b, "%d. %s\n", i+1, line)
	}
	return strings.TrimSpace(b.String())
}
