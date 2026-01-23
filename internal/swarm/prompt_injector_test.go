package swarm

import (
	"context"
	"testing"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/ratelimit"
)

func TestNewPromptInjector(t *testing.T) {
	injector := NewPromptInjector()

	if injector == nil {
		t.Fatal("NewPromptInjector returned nil")
	}

	if injector.TmuxClient != nil {
		t.Error("expected TmuxClient to be nil for default client")
	}

	if injector.StaggerDelay != 300*time.Millisecond {
		t.Errorf("expected StaggerDelay of 300ms, got %v", injector.StaggerDelay)
	}

	if injector.EnterDelay != 100*time.Millisecond {
		t.Errorf("expected EnterDelay of 100ms, got %v", injector.EnterDelay)
	}

	if injector.DoubleEnterDelay != 500*time.Millisecond {
		t.Errorf("expected DoubleEnterDelay of 500ms, got %v", injector.DoubleEnterDelay)
	}

	if injector.Logger == nil {
		t.Error("expected non-nil Logger")
	}

	if len(injector.Templates) == 0 {
		t.Error("expected non-empty Templates map")
	}
}

func TestGetTemplate(t *testing.T) {
	injector := NewPromptInjector()

	tests := []struct {
		name        string
		expectEmpty bool
	}{
		{"default", false},
		{"review", false},
		{"test", false},
		{"nonexistent", false}, // Falls back to default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := injector.GetTemplate(tt.name)
			if tt.expectEmpty && tmpl != "" {
				t.Errorf("expected empty template for %q, got %q", tt.name, tmpl)
			}
			if !tt.expectEmpty && tmpl == "" {
				t.Errorf("expected non-empty template for %q", tt.name)
			}
		})
	}
}

func TestSetTemplate(t *testing.T) {
	injector := NewPromptInjector()

	customTemplate := "This is a custom template"
	injector.SetTemplate("custom", customTemplate)

	result := injector.GetTemplate("custom")
	if result != customTemplate {
		t.Errorf("expected template %q, got %q", customTemplate, result)
	}
}

func TestNeedsDoubleEnter(t *testing.T) {
	tests := []struct {
		agentType string
		expected  bool
	}{
		{"cc", false},
		{"claude", false},
		{"cod", true},
		{"codex", true},
		{"gmi", true},
		{"gemini", true},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.agentType, func(t *testing.T) {
			result := needsDoubleEnter(tt.agentType)
			if result != tt.expected {
				t.Errorf("needsDoubleEnter(%q) = %v, want %v", tt.agentType, result, tt.expected)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	injector := NewPromptInjector()
	result := injector.WithLogger(nil)

	if result != injector {
		t.Error("WithLogger should return the same injector for chaining")
	}
}

func TestWithStaggerDelay(t *testing.T) {
	injector := NewPromptInjector()
	customDelay := 500 * time.Millisecond

	result := injector.WithStaggerDelay(customDelay)

	if result != injector {
		t.Error("WithStaggerDelay should return the same injector for chaining")
	}

	if injector.StaggerDelay != customDelay {
		t.Errorf("expected StaggerDelay of %v, got %v", customDelay, injector.StaggerDelay)
	}
}

func TestInjectionResult(t *testing.T) {
	now := time.Now()
	result := InjectionResult{
		SessionPane: "test:1.5",
		AgentType:   "cc",
		Success:     true,
		Duration:    100 * time.Millisecond,
		SentAt:      now,
	}

	if result.SessionPane != "test:1.5" {
		t.Errorf("unexpected SessionPane: %s", result.SessionPane)
	}
	if result.AgentType != "cc" {
		t.Errorf("unexpected AgentType: %s", result.AgentType)
	}
	if !result.Success {
		t.Error("expected Success to be true")
	}
	if result.Error != "" {
		t.Errorf("expected empty Error, got %q", result.Error)
	}
}

func TestInjectionTarget(t *testing.T) {
	target := InjectionTarget{
		SessionPane: "myproject:1.2",
		AgentType:   "cod",
	}

	if target.SessionPane != "myproject:1.2" {
		t.Errorf("unexpected SessionPane: %s", target.SessionPane)
	}
	if target.AgentType != "cod" {
		t.Errorf("unexpected AgentType: %s", target.AgentType)
	}
}

func TestBatchInjectionResult(t *testing.T) {
	result := BatchInjectionResult{
		TotalPanes: 5,
		Successful: 4,
		Failed:     1,
		Results: []InjectionResult{
			{SessionPane: "s:1.1", AgentType: "cc", Success: true},
			{SessionPane: "s:1.2", AgentType: "cc", Success: true},
			{SessionPane: "s:1.3", AgentType: "cod", Success: true},
			{SessionPane: "s:1.4", AgentType: "cod", Success: true},
			{SessionPane: "s:1.5", AgentType: "gmi", Success: false, Error: "test error"},
		},
		Duration: 2 * time.Second,
	}

	if result.TotalPanes != 5 {
		t.Errorf("expected TotalPanes of 5, got %d", result.TotalPanes)
	}
	if result.Successful != 4 {
		t.Errorf("expected Successful of 4, got %d", result.Successful)
	}
	if result.Failed != 1 {
		t.Errorf("expected Failed of 1, got %d", result.Failed)
	}
	if len(result.Results) != 5 {
		t.Errorf("expected 5 results, got %d", len(result.Results))
	}
}

func TestInjectSwarmNilPlan(t *testing.T) {
	injector := NewPromptInjector()
	result, err := injector.InjectSwarm(nil, "test prompt")

	if err == nil {
		t.Error("expected error for nil plan")
	}
	if result != nil {
		t.Error("expected nil result for nil plan")
	}
}

func TestInjectBatchEmpty(t *testing.T) {
	injector := NewPromptInjector()
	result, err := injector.InjectBatch([]InjectionTarget{}, "test prompt")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.TotalPanes != 0 {
		t.Errorf("expected TotalPanes of 0, got %d", result.TotalPanes)
	}
	if result.Successful != 0 {
		t.Errorf("expected Successful of 0, got %d", result.Successful)
	}
	if result.Failed != 0 {
		t.Errorf("expected Failed of 0, got %d", result.Failed)
	}
}

func TestPromptInjectorTmuxClient(t *testing.T) {
	injector := NewPromptInjector()
	client := injector.tmuxClient()

	if client == nil {
		t.Error("expected non-nil client from tmuxClient()")
	}
}

func TestLoggerHelper(t *testing.T) {
	injector := NewPromptInjector()
	logger := injector.logger()

	if logger == nil {
		t.Error("expected non-nil logger from logger()")
	}
}

func TestDefaultMarchingOrdersNotEmpty(t *testing.T) {
	if DefaultMarchingOrders == "" {
		t.Error("DefaultMarchingOrders should not be empty")
	}

	// Check it contains key instructions
	if len(DefaultMarchingOrders) < 100 {
		t.Error("DefaultMarchingOrders seems too short")
	}
}

func TestPromptTemplateConstants(t *testing.T) {
	if ReviewTemplate == "" {
		t.Error("ReviewTemplate should not be empty")
	}
	if TestTemplate == "" {
		t.Error("TestTemplate should not be empty")
	}
}

func TestWithRateLimitTracker(t *testing.T) {
	injector := NewPromptInjector()
	tracker := ratelimit.NewRateLimitTracker("")

	result := injector.WithRateLimitTracker(tracker)

	if result != injector {
		t.Error("WithRateLimitTracker should return the same injector for chaining")
	}

	if injector.RateLimitTracker != tracker {
		t.Error("expected RateLimitTracker to be set")
	}
}

func TestWithAdaptiveDelay(t *testing.T) {
	injector := NewPromptInjector()

	result := injector.WithAdaptiveDelay(true)

	if result != injector {
		t.Error("WithAdaptiveDelay should return the same injector for chaining")
	}

	if !injector.UseAdaptiveDelay {
		t.Error("expected UseAdaptiveDelay to be true")
	}

	injector.WithAdaptiveDelay(false)
	if injector.UseAdaptiveDelay {
		t.Error("expected UseAdaptiveDelay to be false")
	}
}

func TestGetDelayForAgentFixed(t *testing.T) {
	injector := NewPromptInjector().WithStaggerDelay(500 * time.Millisecond)

	delay := injector.getDelayForAgent("cc")

	if delay != 500*time.Millisecond {
		t.Errorf("expected delay of 500ms, got %v", delay)
	}
}

func TestGetDelayForAgentAdaptive(t *testing.T) {
	tracker := ratelimit.NewRateLimitTracker("")
	injector := NewPromptInjector().
		WithRateLimitTracker(tracker).
		WithAdaptiveDelay(true)

	// Default anthropic delay is 15s
	delay := injector.getDelayForAgent("cc")

	// Should use tracker's optimal delay (defaults to 15s for anthropic)
	if delay != ratelimit.DefaultDelayAnthropic {
		t.Errorf("expected delay of %v for cc, got %v", ratelimit.DefaultDelayAnthropic, delay)
	}

	// Test with openai alias
	delay = injector.getDelayForAgent("cod")
	if delay != ratelimit.DefaultDelayOpenAI {
		t.Errorf("expected delay of %v for cod, got %v", ratelimit.DefaultDelayOpenAI, delay)
	}
}

func TestRecordSuccess(t *testing.T) {
	tracker := ratelimit.NewRateLimitTracker("")
	injector := NewPromptInjector().
		WithRateLimitTracker(tracker).
		WithAdaptiveDelay(true)

	// Record multiple successes
	for i := 0; i < 10; i++ {
		injector.recordSuccess("cc")
	}

	// Check that tracker recorded the successes
	state := tracker.GetProviderState("anthropic")
	if state == nil {
		t.Fatal("expected provider state to exist")
	}

	if state.TotalSuccesses != 10 {
		t.Errorf("expected 10 successes, got %d", state.TotalSuccesses)
	}
}

func TestRecordSuccessNoOpWithoutTracker(t *testing.T) {
	injector := NewPromptInjector().WithAdaptiveDelay(true)

	// Should not panic when tracker is nil
	injector.recordSuccess("cc")
}

func TestRecordSuccessNoOpWithoutAdaptiveDelay(t *testing.T) {
	tracker := ratelimit.NewRateLimitTracker("")
	injector := NewPromptInjector().
		WithRateLimitTracker(tracker).
		WithAdaptiveDelay(false)

	injector.recordSuccess("cc")

	// Should not record anything when adaptive delay is disabled
	state := tracker.GetProviderState("anthropic")
	if state != nil {
		t.Error("expected no provider state when adaptive delay is disabled")
	}
}

func TestInjectBatchWithContextCancellation(t *testing.T) {
	injector := NewPromptInjector()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	targets := []InjectionTarget{
		{SessionPane: "test:1.1", AgentType: "cc"},
		{SessionPane: "test:1.2", AgentType: "cc"},
	}

	result, err := injector.InjectBatchWithContext(ctx, targets, "test prompt")

	if err == nil {
		t.Error("expected error for cancelled context")
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestInjectSwarmWithContext(t *testing.T) {
	injector := NewPromptInjector()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	plan := &SwarmPlan{
		Sessions: []SessionSpec{
			{
				Name:      "test",
				AgentType: "cc",
				PaneCount: 2,
				Panes: []PaneSpec{
					{Index: 1, AgentType: "cc"},
					{Index: 2, AgentType: "cc"},
				},
			},
		},
	}

	result, err := injector.InjectSwarmWithContext(ctx, plan, "test prompt")

	if err == nil {
		t.Error("expected error for cancelled context")
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
