package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Dicklesworthstone/ntm/internal/tmux"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var panes int

	cmd := &cobra.Command{
		Use:   "create <session-name>",
		Short: "Create a new tmux session with multiple panes",
		Long: `Create a new tmux session with the specified number of panes.
The session directory is created under PROJECTS_BASE if it doesn't exist.

Example:
  ntm create myproject           # Create with default panes
  ntm create myproject --panes=6 # Create with 6 panes`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args[0], panes)
		},
	}

	cmd.Flags().IntVarP(&panes, "panes", "p", 0, "number of panes to create (default from config)")

	return cmd
}

func runCreate(session string, panes int) error {
	if err := tmux.EnsureInstalled(); err != nil {
		return err
	}

	if err := tmux.ValidateSessionName(session); err != nil {
		return err
	}

	// Get pane count from config if not specified
	if panes <= 0 {
		panes = cfg.Tmux.DefaultPanes
	}

	dir := cfg.GetProjectDir(session)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Directory not found: %s\n", dir)
		if !confirm("Create it?") {
			fmt.Println("Aborted.")
			return nil
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		fmt.Printf("Created %s\n", dir)
	}

	// Check if session already exists
	if tmux.SessionExists(session) {
		fmt.Printf("Session '%s' already exists\n", session)
		return tmux.AttachOrSwitch(session)
	}

	fmt.Printf("Creating session '%s' with %d pane(s)...\n", session, panes)

	// Create the session
	if err := tmux.CreateSession(session, dir); err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	// Add additional panes
	if panes > 1 {
		for i := 1; i < panes; i++ {
			if _, err := tmux.SplitWindow(session, dir); err != nil {
				return fmt.Errorf("creating pane %d: %w", i+1, err)
			}
		}
	}

	fmt.Printf("Created session '%s' with %d pane(s)\n", session, panes)
	return tmux.AttachOrSwitch(session)
}

// confirm prompts the user for y/n confirmation
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
