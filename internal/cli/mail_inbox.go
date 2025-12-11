package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/agentmail"
	"github.com/Dicklesworthstone/ntm/internal/output"
	"github.com/Dicklesworthstone/ntm/internal/tmux"
	"github.com/spf13/cobra"
)

func newMailInboxCmd() *cobra.Command {
	var (
		sessionAgents bool
		agentFilter   string
		urgent        bool
		limit         int
	)

	cmd := &cobra.Command{
		Use:   "inbox [session]",
		Short: "Show aggregate project inbox",
		Long: `Display an aggregate inbox view showing messages across all agents in the project.

		This provides visibility into agent-to-agent communication and Human Overseer messages.
		Messages are deduplicated (shown once even if sent to multiple agents).

		Examples:
		  ntm mail inbox myproject
		  ntm mail inbox myproject --session-agents  # Only show messages for agents in this session
		  ntm mail inbox myproject --agent BlueLake  # Only show messages involving BlueLake
		  ntm mail inbox --urgent                    # Only show urgent messages`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var session string
			if len(args) > 0 {
				session = args[0]
			}
			return runMailInbox(cmd, session, sessionAgents, agentFilter, urgent, limit)
		},
	}

	cmd.Flags().BoolVar(&sessionAgents, "session-agents", false, "filter to agents present in the session")
	cmd.Flags().StringVar(&agentFilter, "agent", "", "filter to messages involving specific agent")
	cmd.Flags().BoolVar(&urgent, "urgent", false, "show only urgent messages")
	cmd.Flags().IntVar(&limit, "limit", 50, "max messages per agent inbox to fetch")

	return cmd
}

type aggregatedMessage struct {
	ID          int
	Subject     string
	From        string
	CreatedTS   time.Time
	Importance  string
	AckRequired bool
	Kind        string
	BodyMD      string // truncated for display
	Recipients  []string
}

func runMailInbox(cmd *cobra.Command, session string, sessionAgents bool, agentFilter string, urgent bool, limit int) error {
	projectKey, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	client := agentmail.NewClient(agentmail.WithProjectKey(projectKey))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if !client.IsAvailable() {
		return fmt.Errorf("agent mail server not available")
	}

	// 1. Discover relevant agents
	allAgents, err := client.ListProjectAgents(ctx, projectKey)
	if err != nil {
		return fmt.Errorf("listing project agents: %w", err)
	}

	var targetAgents []string
	sessionAgentSet := make(map[string]bool)

	// If filtering by session agents, resolve session panes
	if sessionAgents {
		if session == "" {
			if tmux.InTmux() {
				session = tmux.GetCurrentSession()
			} else {
				return fmt.Errorf("session name required for --session-agents (or run inside tmux)")
			}
		}
		panes, err := tmux.GetPanes(session)
		if err != nil {
			return fmt.Errorf("getting session panes: %w", err)
		}
		for _, p := range panes {
			name := resolveAgentName(p)
			if name != "" {
				sessionAgentSet[name] = true
			}
		}
		if len(sessionAgentSet) == 0 {
			return fmt.Errorf("no agents found in session '%s'", session)
		}
	}

	for _, a := range allAgents {
		// "HumanOverseer" doesn't have an inbox we care about usually, but let's include all regular agents
		if a.Name == "HumanOverseer" {
			continue
		}
		// Apply filters for FETCHING
		// If --agent is set, we only need to fetch THAT agent's inbox?
		// No, requirement says "Aggregate View... showing messages... across ALL agents".
		// BUT "Filter to specific agent" might mean "Show messages where X is sender OR recipient".
		// If I only fetch X's inbox, I see messages TO X. I miss messages FROM X to Y.
		// To see messages FROM X, I need Y's inbox.
		// So to do it right, we should fetch ALL inboxes and then filter.
		// Optimization: if filter is specific agent, can we optimize? Maybe not easily without search.
		// Let's stick to fetching all inboxes (limit is key).
		targetAgents = append(targetAgents, a.Name)
	}

	// 2. Fetch inboxes in parallel
	var (
		mu       sync.Mutex
		messages = make(map[int]*aggregatedMessage)
		wg       sync.WaitGroup
	)

	errChan := make(chan error, len(targetAgents))

	for _, agentName := range targetAgents {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			inbox, err := client.FetchInbox(ctx, agentmail.FetchInboxOptions{
				ProjectKey: projectKey,
				AgentName:  name,
				UrgentOnly: urgent,
				Limit:      limit,
			})
			if err != nil {
				// Log warning but don't fail completely?
				// errChan <- fmt.Errorf("fetching inbox for %s: %w", name, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			for _, msg := range inbox {
				agg, exists := messages[msg.ID]
				if !exists {
					agg = &aggregatedMessage{
						ID:          msg.ID,
						Subject:     msg.Subject,
						From:        msg.From,
						CreatedTS:   msg.CreatedTS,
						Importance:  msg.Importance,
						AckRequired: msg.AckRequired,
						Kind:        msg.Kind,
						BodyMD:      msg.BodyMD, // Empty unless we ask for it
						Recipients:  []string{},
					}
					messages[msg.ID] = agg
				}
				// Add this agent as a recipient
				// Check if already in list (unlikely unless agent listed twice)
				isPresent := false
				for _, r := range agg.Recipients {
					if r == name {
						isPresent = true
						break
					}
				}
				if !isPresent {
					agg.Recipients = append(agg.Recipients, name)
				}
			}
		}(agentName)
	}

	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		// report first error? Or just continue with partial data?
		// For CLI, partial data is better than nothing usually.
	}

	// 3. Flatten and Filter
	var result []*aggregatedMessage
	for _, m := range messages {
		// Apply filters
		if urgent && m.Importance != "urgent" && m.Importance != "high" {
			continue
		}

		if agentFilter != "" {
			// Check if sender matches
			senderMatch := strings.EqualFold(m.From, agentFilter)
			// Check if any recipient matches
			recipientMatch := false
			for _, r := range m.Recipients {
				if strings.EqualFold(r, agentFilter) {
					recipientMatch = true
					break
				}
			}
			if !senderMatch && !recipientMatch {
				continue
			}
		}

		if sessionAgents {
			// Check if ANY recipient or sender is in the session
			relevant := false
			if sessionAgentSet[m.From] {
				relevant = true
			} else {
				for _, r := range m.Recipients {
					if sessionAgentSet[r] {
						relevant = true
						break
					}
				}
			}
			if !relevant {
				continue
			}
		}

		result = append(result, m)
	}

	// 4. Sort (Newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedTS.After(result[j].CreatedTS)
	})

	// 5. Output
	if IsJSONOutput() {
		return encodeJSONResult(mailJSONWriter(cmd), result)
	}

	if len(result) == 0 {
		fmt.Println("Inbox empty.")
		return nil
	}

	// TUI Display
	fmt.Printf("┌─────────────────────────────────────────────────────────────┐\n")
	title := fmt.Sprintf("Project Inbox: %s (%d agents)", filepath.Base(projectKey), len(allAgents))
	fmt.Printf("│ %-59s │\n", title)
	
	for _, m := range result {
		fmt.Printf("├─────────────────────────────────────────────────────────────┤\n")
		
		// Status icon
		icon := "○" // Read/Normal
		if m.Importance == "urgent" || m.Importance == "high" {
			icon = "●" // Urgent
		}
		// Ack check? We don't know if *this* viewer acked it, but we know if ack required.
		
		urgencyTag := ""
		if m.Importance == "urgent" || m.Importance == "high" {
			urgencyTag = fmt.Sprintf("[%s] ", strings.ToUpper(m.Importance))
		}

		// Line 1: Icon ID [URGENT] Subject
		// Truncate subject to fit
		sub := m.Subject
		if len(sub) > 40 {
			sub = sub[:37] + "..."
		}
		fmt.Printf("│ %s #%d %s%s\n", icon, m.ID, urgencyTag, sub)

		// Line 2: From -> To
		recipientsStr := strings.Join(m.Recipients, ", ")
		if len(recipientsStr) > 40 {
			recipientsStr = recipientsStr[:37] + "..."
		}
		fmt.Printf("│   %s → %s\n", m.From, recipientsStr)

		// Line 3: Time | Thread | Flags
		timeStr := m.CreatedTS.Format("15:04")
		ago := time.Since(m.CreatedTS).Round(time.Minute)
		if ago < time.Hour {
			timeStr = fmt.Sprintf("%dm ago", int(ago.Minutes()))
		}
		
		ackStr := ""
		if m.AckRequired {
			ackStr = " │ ⚠️ Ack required"
		}
		
		// Note: Thread ID not available in InboxMessage yet, omitted
		fmt.Printf("│   %s%s\n", timeStr, ackStr)
	}
	fmt.Printf("└─────────────────────────────────────────────────────────────┘\n")
	
	// Show active filter info
	filterInfo := []string{}
	if agentFilter != "" {
		filterInfo = append(filterInfo, "Filter: "+agentFilter)
	}
	if sessionAgents {
		filterInfo = append(filterInfo, "Filter: Session Agents")
	}
	if urgent {
		filterInfo = append(filterInfo, "Filter: Urgent")
	}
	if len(filterInfo) > 0 {
		fmt.Printf("Showing messages: %s\n", strings.Join(filterInfo, ", "))
	}

	return nil
}

func resolveAgentName(p tmux.Pane) string {
	if p.Variant != "" {
		return p.Variant
	}
	return ""
}

func encodeJSONResult(w io.Writer, v interface{}) error {
	return output.WriteJSON(w, v, true)
}

func mailJSONWriter(cmd *cobra.Command) io.Writer {
	return cmd.OutOrStdout()
}
