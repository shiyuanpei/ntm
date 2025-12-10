package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/tmux"
	"github.com/Dicklesworthstone/ntm/tests/testutil"
)

// TestCLISpawnSendAndStatus verifies that ntm CLI commands drive tmux correctly:
// - spawn creates a tmux session
// - we can add synthetic agent panes
// - send targets those agent panes
// - status reports the expected pane count.
func TestCLISpawnSendAndStatus(t *testing.T) {
	testutil.RequireNTMBinary(t)
	testutil.RequireTmux(t)

	logger := testutil.NewTestLogger(t, t.TempDir())

	// Create a session without launching external agent binaries.
	session := testutil.CreateTestSession(t, logger, testutil.SessionConfig{})
	testutil.AssertSessionExists(t, logger, session)

	// Add synthetic cc and cod panes so send/status have targets.
	ccPane, err := tmux.SplitWindow(session, t.TempDir())
	if err != nil {
		t.Fatalf("failed to split window for cc pane: %v", err)
	}
	if err := tmux.SetPaneTitle(ccPane, fmt.Sprintf("%s__cc_1", session)); err != nil {
		t.Fatalf("failed to set cc pane title: %v", err)
	}

	codPane, err := tmux.SplitWindow(session, t.TempDir())
	if err != nil {
		t.Fatalf("failed to split window for cod pane: %v", err)
	}
	if err := tmux.SetPaneTitle(codPane, fmt.Sprintf("%s__cod_1", session)); err != nil {
		t.Fatalf("failed to set cod pane title: %v", err)
	}

	_ = tmux.ApplyTiledLayout(session)

	// status should see user + cc + cod = 3 panes.
	testutil.AssertNTMStatus(t, logger, session, 3)

	// Send a command to cc panes and verify it lands.
	const marker = "INTEGRATION_CC_OK"
	testutil.AssertCommandSuccess(t, logger, "ntm", "send", session, "--cc", "echo "+marker)

	testutil.AssertEventually(t, logger, 5*time.Second, 150*time.Millisecond, "cc pane receives send payload", func() bool {
		out, err := tmux.CapturePaneOutput(ccPane, 200)
		if err != nil {
			return false
		}
		return strings.Contains(out, marker)
	})
}
