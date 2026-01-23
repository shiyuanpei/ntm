package swarm

import (
	"context"
	"testing"
	"time"
)

func TestNewLimitDetector(t *testing.T) {
	detector := NewLimitDetector()

	if detector == nil {
		t.Fatal("NewLimitDetector returned nil")
	}

	if detector.TmuxClient != nil {
		t.Error("expected TmuxClient to be nil for default client")
	}

	if detector.CheckInterval != 5*time.Second {
		t.Errorf("expected CheckInterval of 5s, got %v", detector.CheckInterval)
	}

	if detector.CaptureLines != 50 {
		t.Errorf("expected CaptureLines of 50, got %d", detector.CaptureLines)
	}

	if detector.eventChan == nil {
		t.Error("expected eventChan to be initialized")
	}

	if detector.monitoredPanes == nil {
		t.Error("expected monitoredPanes to be initialized")
	}
}

func TestLimitDetectorEvents(t *testing.T) {
	detector := NewLimitDetector()
	eventChan := detector.Events()

	if eventChan == nil {
		t.Error("Events() returned nil channel")
	}
}

func TestLimitDetectorStartNilPlan(t *testing.T) {
	detector := NewLimitDetector()
	ctx := context.Background()

	err := detector.Start(ctx, nil)
	if err != nil {
		t.Errorf("Start with nil plan should not error, got: %v", err)
	}
}

func TestLimitDetectorStartEmptyPlan(t *testing.T) {
	detector := NewLimitDetector()
	ctx := context.Background()

	plan := &SwarmPlan{
		Sessions: []SessionSpec{},
	}

	err := detector.Start(ctx, plan)
	if err != nil {
		t.Errorf("Start with empty plan should not error, got: %v", err)
	}

	// Should have no monitored panes
	panes := detector.MonitoredPanes()
	if len(panes) != 0 {
		t.Errorf("expected 0 monitored panes, got %d", len(panes))
	}
}

func TestLimitDetectorStop(t *testing.T) {
	detector := NewLimitDetector()
	ctx := context.Background()

	plan := &SwarmPlan{
		Sessions: []SessionSpec{
			{
				Name:      "test_session",
				AgentType: "cc",
				Panes: []PaneSpec{
					{Index: 1, AgentType: "cc"},
				},
			},
		},
	}

	// Start and then stop
	_ = detector.Start(ctx, plan)
	detector.Stop()

	// Should have no monitored panes after stop
	panes := detector.MonitoredPanes()
	if len(panes) != 0 {
		t.Errorf("expected 0 monitored panes after Stop, got %d", len(panes))
	}
}

func TestLimitDetectorIsMonitoring(t *testing.T) {
	detector := NewLimitDetector()

	// Should not be monitoring anything initially
	if detector.IsMonitoring("test:1.1") {
		t.Error("expected IsMonitoring to return false for unmonitored pane")
	}
}

func TestLimitDetectorMonitoredPanes(t *testing.T) {
	detector := NewLimitDetector()

	panes := detector.MonitoredPanes()
	if panes == nil {
		t.Error("MonitoredPanes() returned nil")
	}
	if len(panes) != 0 {
		t.Errorf("expected 0 monitored panes initially, got %d", len(panes))
	}
}

func TestLimitDetectorStopPane(t *testing.T) {
	detector := NewLimitDetector()

	// StopPane on non-existent pane should not panic
	detector.StopPane("nonexistent:1.1")
}

func TestLimitEvent(t *testing.T) {
	event := LimitEvent{
		SessionPane: "test:1.5",
		AgentType:   "cc",
		Pattern:     "rate limit",
		RawOutput:   "You've hit your rate limit",
		DetectedAt:  time.Now(),
	}

	if event.SessionPane != "test:1.5" {
		t.Errorf("unexpected SessionPane: %s", event.SessionPane)
	}
	if event.AgentType != "cc" {
		t.Errorf("unexpected AgentType: %s", event.AgentType)
	}
	if event.Pattern != "rate limit" {
		t.Errorf("unexpected Pattern: %s", event.Pattern)
	}
}

func TestGetPatternsForAgent(t *testing.T) {
	detector := NewLimitDetector()

	tests := []struct {
		agentType     string
		expectDefault bool
	}{
		{"cc", false},
		{"cod", false},
		{"gmi", false},
		{"claude", false},
		{"codex", false},
		{"gemini", false},
		{"unknown", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.agentType, func(t *testing.T) {
			patterns := detector.getPatternsForAgent(tt.agentType)
			if len(patterns) == 0 {
				t.Error("expected non-empty patterns")
			}
		})
	}
}

func TestCheckOutputEmpty(t *testing.T) {
	detector := NewLimitDetector()

	event := detector.checkOutput("test:1.1", "cc", "")
	if event != nil {
		t.Error("expected nil event for empty output")
	}
}

func TestCheckOutputNoMatch(t *testing.T) {
	detector := NewLimitDetector()

	output := "Normal agent output\nNo issues here\nJust working on code"
	event := detector.checkOutput("test:1.1", "cc", output)
	if event != nil {
		t.Error("expected nil event for output with no rate limit patterns")
	}
}

func TestCheckOutputMatch(t *testing.T) {
	detector := NewLimitDetector()

	output := "Working on task...\nYou've hit your rate limit. Please wait.\nTry again later."
	event := detector.checkOutput("test:1.1", "cc", output)

	if event == nil {
		t.Fatal("expected non-nil event for output with rate limit pattern")
	}

	if event.SessionPane != "test:1.1" {
		t.Errorf("unexpected SessionPane: %s", event.SessionPane)
	}
	if event.AgentType != "cc" {
		t.Errorf("unexpected AgentType: %s", event.AgentType)
	}
	if event.RawOutput != output {
		t.Error("expected RawOutput to match input")
	}
}

func TestCheckOutputCaseInsensitive(t *testing.T) {
	detector := NewLimitDetector()

	// Test case insensitivity
	output := "RATE LIMIT EXCEEDED"
	event := detector.checkOutput("test:1.1", "cc", output)

	if event == nil {
		t.Error("expected pattern matching to be case insensitive")
	}
}

func TestDefaultLimitPatterns(t *testing.T) {
	if len(defaultLimitPatterns) == 0 {
		t.Error("expected non-empty defaultLimitPatterns")
	}

	expectedPatterns := []string{
		"rate limit",
		"usage limit",
		"quota exceeded",
		"too many requests",
	}

	for _, expected := range expectedPatterns {
		found := false
		for _, pattern := range defaultLimitPatterns {
			if pattern == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected pattern %q in defaultLimitPatterns", expected)
		}
	}
}

func TestLimitDetectorTmuxClientHelper(t *testing.T) {
	// With nil client, should return default
	detector := NewLimitDetector()
	client := detector.tmuxClient()
	if client == nil {
		t.Error("expected non-nil client from tmuxClient()")
	}
}

func TestLimitDetectorLoggerHelper(t *testing.T) {
	// With nil logger, should return default
	detector := &LimitDetector{}
	logger := detector.logger()
	if logger == nil {
		t.Error("expected non-nil logger from logger()")
	}
}
