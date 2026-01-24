package ensemble

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/Dicklesworthstone/ntm/internal/tmux"
)

// EarlyStopConfig controls optional early-stopping behavior for ensembles.
// This is a placeholder for future policy logic.
type EarlyStopConfig struct {
	Enabled     bool    `json:"enabled" toml:"enabled" yaml:"enabled"`
	MinUtility  float64 `json:"min_utility,omitempty" toml:"min_utility" yaml:"min_utility,omitempty"`
	MaxStagnant int     `json:"max_stagnant,omitempty" toml:"max_stagnant" yaml:"max_stagnant,omitempty"`
}

// EnsembleConfig defines the inputs for starting an ensemble run.
type EnsembleConfig struct {
	SessionName string
	Question    string
	Ensemble    string   // built-in or user-defined ensemble name
	Modes       []string // explicit mode IDs or codes (if not using ensemble)

	AgentMix   map[string]int
	Assignment string // round-robin, affinity, explicit

	Synthesis SynthesisConfig
	Budget    BudgetConfig
	Cache     CacheConfig
	EarlyStop EarlyStopConfig
}

// EnsembleManager orchestrates ensemble sessions: spawn, assign, inject, persist.
type EnsembleManager struct {
	Catalog          *ModeCatalog
	Registry         *EnsembleRegistry
	TmuxClient       *tmux.Client
	Logger           *slog.Logger
	PreambleEngine   *PreambleEngine
	LaunchDelay      time.Duration
	EnterDelay       time.Duration
	DoubleEnterDelay time.Duration
}

// NewEnsembleManager creates a manager with default dependencies.
func NewEnsembleManager(catalog *ModeCatalog, registry *EnsembleRegistry) *EnsembleManager {
	if catalog == nil {
		catalog, _ = GlobalCatalog()
	}
	if registry == nil && catalog != nil {
		registry, _ = GlobalEnsembleRegistry()
	}
	return &EnsembleManager{
		Catalog:          catalog,
		Registry:         registry,
		TmuxClient:       nil,
		Logger:           slog.Default(),
		PreambleEngine:   NewPreambleEngine(),
		LaunchDelay:      200 * time.Millisecond,
		EnterDelay:       100 * time.Millisecond,
		DoubleEnterDelay: 500 * time.Millisecond,
	}
}

func (m *EnsembleManager) tmuxClient() *tmux.Client {
	if m != nil && m.TmuxClient != nil {
		return m.TmuxClient
	}
	return tmux.DefaultClient
}

func (m *EnsembleManager) logger() *slog.Logger {
	if m != nil && m.Logger != nil {
		return m.Logger
	}
	return slog.Default()
}

// SpawnEnsemble creates a tmux session, assigns modes, injects prompts, and persists state.
// Returns the ensemble session state and any error encountered (partial success possible).
func (m *EnsembleManager) SpawnEnsemble(cfg EnsembleConfig) (*EnsembleSession, error) {
	if err := m.validateConfig(cfg); err != nil {
		return nil, err
	}

	catalog, err := m.ensureCatalog()
	if err != nil {
		return nil, err
	}

	registry, err := m.ensureRegistry(catalog)
	if err != nil {
		return nil, err
	}

	modeIDs, presetName, err := resolveModeIDs(cfg, registry, catalog)
	if err != nil {
		return nil, err
	}

	budget := mergeBudgetDefaults(cfg.Budget, DefaultBudgetConfig())
	synthesis := mergeSynthesisDefaults(cfg.Synthesis, DefaultSynthesisConfig())

	agentMix, err := normalizeAgentMix(cfg.AgentMix, len(modeIDs))
	if err != nil {
		return nil, err
	}

	client := m.tmuxClient()
	if client.SessionExists(cfg.SessionName) {
		return nil, fmt.Errorf("session %q already exists", cfg.SessionName)
	}

	workDir := currentDir()
	if workDir == "" {
		workDir = "/tmp"
	}

	if err := client.CreateSession(cfg.SessionName, workDir); err != nil {
		return nil, fmt.Errorf("create session %q: %w", cfg.SessionName, err)
	}

	m.logger().Info("ensemble session created",
		"session", cfg.SessionName,
		"modes", len(modeIDs),
		"agent_mix", agentMix)

	paneErrs := m.spawnAgentPanes(cfg.SessionName, workDir, agentMix)

	panes, err := client.GetPanes(cfg.SessionName)
	if err != nil {
		return nil, fmt.Errorf("get panes for session %q: %w", cfg.SessionName, err)
	}

	assignments, err := assignModes(cfg, modeIDs, panes, catalog)
	if err != nil {
		return nil, err
	}

	state := &EnsembleSession{
		SessionName:       cfg.SessionName,
		Question:          cfg.Question,
		PresetUsed:        presetName,
		Assignments:       assignments,
		Status:            EnsembleInjecting,
		SynthesisStrategy: synthesis.Strategy,
		CreatedAt:         time.Now().UTC(),
	}

	injectErrs := m.injectAssignments(state, catalog, cfg.Question, budget.MaxTokensPerMode)

	activeCount := 0
	for _, assignment := range state.Assignments {
		if assignment.Status == AssignmentActive {
			activeCount++
		}
	}

	if activeCount == 0 {
		state.Status = EnsembleError
	} else {
		state.Status = EnsembleActive
	}

	allErrs := append([]error{}, paneErrs...)
	allErrs = append(allErrs, injectErrs...)
	if len(allErrs) > 0 {
		state.Error = errors.Join(allErrs...).Error()
	}

	if err := SaveSession(cfg.SessionName, state); err != nil {
		return state, err
	}

	if len(allErrs) > 0 {
		return state, errors.Join(allErrs...)
	}
	return state, nil
}

func (m *EnsembleManager) validateConfig(cfg EnsembleConfig) error {
	if strings.TrimSpace(cfg.SessionName) == "" {
		return errors.New("session_name is required")
	}
	if err := tmux.ValidateSessionName(cfg.SessionName); err != nil {
		return err
	}
	if strings.TrimSpace(cfg.Question) == "" {
		return errors.New("question is required")
	}
	if cfg.Ensemble == "" && len(cfg.Modes) == 0 {
		return errors.New("either ensemble or modes must be provided")
	}
	if cfg.Ensemble != "" && len(cfg.Modes) > 0 {
		return errors.New("provide either ensemble or modes, not both")
	}
	if cfg.Synthesis.Strategy != "" && !cfg.Synthesis.Strategy.IsValid() {
		return fmt.Errorf("invalid synthesis strategy %q", cfg.Synthesis.Strategy)
	}
	if cfg.Budget.MaxTokensPerMode < 0 || cfg.Budget.MaxTotalTokens < 0 || cfg.Budget.MaxRetries < 0 {
		return errors.New("budget values cannot be negative")
	}

	strategy := strings.ToLower(strings.TrimSpace(cfg.Assignment))
	switch strategy {
	case "", "round-robin", "roundrobin", "affinity", "category", "explicit":
		if strategy == "explicit" && len(cfg.Modes) == 0 {
			return errors.New("explicit assignment requires mode:agent mappings")
		}
		return nil
	default:
		return fmt.Errorf("unsupported assignment strategy %q", cfg.Assignment)
	}
}

func (m *EnsembleManager) ensureCatalog() (*ModeCatalog, error) {
	if m != nil && m.Catalog != nil {
		return m.Catalog, nil
	}
	return GlobalCatalog()
}

func (m *EnsembleManager) ensureRegistry(catalog *ModeCatalog) (*EnsembleRegistry, error) {
	if m != nil && m.Registry != nil {
		return m.Registry, nil
	}
	return GlobalEnsembleRegistry()
}

func resolveModeIDs(cfg EnsembleConfig, registry *EnsembleRegistry, catalog *ModeCatalog) ([]string, string, error) {
	if cfg.Ensemble != "" {
		if registry == nil {
			return nil, "", errors.New("ensemble registry not available")
		}
		preset := registry.Get(cfg.Ensemble)
		if preset == nil {
			return nil, "", fmt.Errorf("ensemble %q not found", cfg.Ensemble)
		}
		ids, err := preset.ResolveIDs(catalog)
		if err != nil {
			return nil, "", err
		}
		return ids, preset.Name, nil
	}

	assignment := strings.ToLower(strings.TrimSpace(cfg.Assignment))
	if assignment == "explicit" {
		ids, err := resolveExplicitModeIDs(cfg.Modes, catalog)
		return ids, "", err
	}

	modeIDs := make([]string, 0, len(cfg.Modes))
	seen := make(map[string]bool, len(cfg.Modes))
	for _, raw := range cfg.Modes {
		modeID, _, err := resolveMode(raw, catalog)
		if err != nil {
			return nil, "", err
		}
		if seen[modeID] {
			return nil, "", fmt.Errorf("duplicate mode %q", modeID)
		}
		seen[modeID] = true
		modeIDs = append(modeIDs, modeID)
	}
	sort.Strings(modeIDs)
	return modeIDs, "", nil
}

func resolveExplicitModeIDs(specs []string, catalog *ModeCatalog) ([]string, error) {
	expanded := expandSpecs(specs)
	if len(expanded) == 0 {
		return nil, errors.New("explicit assignment requires at least one mapping")
	}
	ids := make([]string, 0, len(expanded))
	seen := make(map[string]bool, len(expanded))
	for _, spec := range expanded {
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid assignment %q: expected mode:agent", spec)
		}
		modeID, _, err := resolveMode(parts[0], catalog)
		if err != nil {
			return nil, err
		}
		if seen[modeID] {
			return nil, fmt.Errorf("duplicate mode %q", modeID)
		}
		seen[modeID] = true
		ids = append(ids, modeID)
	}
	sort.Strings(ids)
	return ids, nil
}

func normalizeAgentMix(agentMix map[string]int, modeCount int) (map[string]int, error) {
	if len(agentMix) == 0 {
		return map[string]int{
			string(tmux.AgentClaude): modeCount,
		}, nil
	}

	normalized := make(map[string]int, len(agentMix))
	total := 0
	for rawType, count := range agentMix {
		if count < 0 {
			return nil, fmt.Errorf("agent mix for %q cannot be negative", rawType)
		}
		if count == 0 {
			continue
		}
		agentType := strings.ToLower(strings.TrimSpace(rawType))
		if agentType == "" {
			return nil, errors.New("agent type cannot be empty")
		}
		normalized[agentType] += count
		total += count
	}

	if total == 0 {
		return nil, errors.New("agent mix must include at least one agent")
	}
	if total < modeCount {
		return nil, fmt.Errorf("agent mix has %d panes but %d modes requested", total, modeCount)
	}

	return normalized, nil
}

func assignModes(cfg EnsembleConfig, modeIDs []string, panes []tmux.Pane, catalog *ModeCatalog) ([]ModeAssignment, error) {
	strategy := strings.ToLower(strings.TrimSpace(cfg.Assignment))
	switch strategy {
	case "", "round-robin", "roundrobin":
		assignments := AssignRoundRobin(modeIDs, panes)
		if assignments == nil {
			return nil, errors.New("round-robin assignment failed")
		}
		return assignments, nil
	case "affinity", "category":
		assignments := AssignByCategory(modeIDs, panes, catalog)
		if assignments == nil {
			return nil, errors.New("category assignment failed")
		}
		return assignments, nil
	case "explicit":
		canonicalSpecs, err := canonicalExplicitSpecs(cfg.Modes, catalog)
		if err != nil {
			return nil, err
		}
		assignments, err := AssignExplicit(canonicalSpecs, panes)
		if err != nil {
			return nil, err
		}
		return assignments, nil
	default:
		return nil, fmt.Errorf("unsupported assignment strategy %q", cfg.Assignment)
	}
}

func canonicalExplicitSpecs(specs []string, catalog *ModeCatalog) ([]string, error) {
	expanded := expandSpecs(specs)
	if len(expanded) == 0 {
		return nil, errors.New("explicit assignment requires at least one mapping")
	}
	canonical := make([]string, 0, len(expanded))
	for _, spec := range expanded {
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid assignment %q: expected mode:agent", spec)
		}
		modeID, _, err := resolveMode(parts[0], catalog)
		if err != nil {
			return nil, err
		}
		agentType := strings.ToLower(strings.TrimSpace(parts[1]))
		if agentType == "" {
			return nil, fmt.Errorf("invalid assignment %q: empty agent type", spec)
		}
		canonical = append(canonical, fmt.Sprintf("%s:%s", modeID, agentType))
	}
	return canonical, nil
}

func (m *EnsembleManager) spawnAgentPanes(sessionName, workDir string, agentMix map[string]int) []error {
	client := m.tmuxClient()
	order := orderedAgentTypes(agentMix)
	var errs []error

	panes, err := client.GetPanes(sessionName)
	if err != nil || len(panes) == 0 {
		err = fmt.Errorf("get initial pane for session %q: %w", sessionName, err)
		return []error{err}
	}

	firstPaneID := panes[0].ID
	usedFirst := false
	indices := make(map[string]int)

	for _, agentType := range order {
		count := agentMix[agentType]
		for i := 0; i < count; i++ {
			var paneID string
			if !usedFirst {
				paneID = firstPaneID
				usedFirst = true
			} else {
				paneID, err = client.SplitWindow(sessionName, workDir)
				if err != nil {
					errs = append(errs, fmt.Errorf("split window for %s: %w", agentType, err))
					continue
				}
			}

			indices[agentType]++
			title := tmux.FormatPaneName(sessionName, agentType, indices[agentType], "")
			if err := client.SetPaneTitle(paneID, title); err != nil {
				errs = append(errs, fmt.Errorf("set pane title %s: %w", paneID, err))
			}

			if err := m.launchAgentInPane(paneID, agentType); err != nil {
				errs = append(errs, fmt.Errorf("launch agent %s: %w", paneID, err))
			}

			if m.LaunchDelay > 0 {
				time.Sleep(m.LaunchDelay)
			}
		}
	}

	_ = client.ApplyTiledLayout(sessionName)
	return errs
}

func (m *EnsembleManager) launchAgentInPane(paneID, agentType string) error {
	client := m.tmuxClient()
	if err := client.SendKeys(paneID, agentType, false); err != nil {
		return err
	}
	if m.EnterDelay > 0 {
		time.Sleep(m.EnterDelay)
	}
	if err := client.SendKeys(paneID, "", true); err != nil {
		return err
	}
	return nil
}

func (m *EnsembleManager) injectAssignments(state *EnsembleSession, catalog *ModeCatalog, question string, tokenCap int) []error {
	client := m.tmuxClient()
	panes, err := client.GetPanes(state.SessionName)
	if err != nil {
		return []error{fmt.Errorf("get panes for session %q: %w", state.SessionName, err)}
	}

	paneIDs := make(map[string]string, len(panes)*2)
	for _, pane := range panes {
		if pane.Title != "" {
			paneIDs[pane.Title] = pane.ID
		}
		if pane.ID != "" {
			paneIDs[pane.ID] = pane.ID
		}
	}

	var errs []error

	for i := range state.Assignments {
		assignment := &state.Assignments[i]
		mode := catalog.GetMode(assignment.ModeID)
		if mode == nil {
			err := fmt.Errorf("mode %q not found", assignment.ModeID)
			assignment.Status = AssignmentError
			assignment.Error = err.Error()
			errs = append(errs, err)
			continue
		}

		assignment.Status = AssignmentInjecting
		preamble, err := m.PreambleEngine.Render(&PreambleData{
			Problem:      question,
			ContextPack:  nil,
			Mode:         mode,
			TokenCap:     tokenCap,
			OutputSchema: GetSchemaContract(),
		})
		if err != nil {
			assignment.Status = AssignmentError
			assignment.Error = err.Error()
			errs = append(errs, err)
			m.logger().Error("ensemble preamble render failed",
				"mode_id", assignment.ModeID,
				"pane", assignment.PaneName,
				"error", err)
			continue
		}

		prompt := preamble
		if strings.TrimSpace(question) != "" {
			prompt = preamble + "\n\n" + question
		}

		target := paneIDs[assignment.PaneName]
		if target == "" {
			target = assignment.PaneName
		}

		if err := sendPrompt(client, target, assignment.AgentType, prompt, m.EnterDelay, m.DoubleEnterDelay); err != nil {
			assignment.Status = AssignmentError
			assignment.Error = err.Error()
			errs = append(errs, err)
			m.logger().Error("ensemble prompt injection failed",
				"mode_id", assignment.ModeID,
				"pane", assignment.PaneName,
				"error", err)
			continue
		}

		assignment.Status = AssignmentActive
		m.logger().Info("ensemble prompt injected",
			"mode_id", assignment.ModeID,
			"pane", assignment.PaneName,
			"agent_type", assignment.AgentType)
	}

	return errs
}

func orderedAgentTypes(agentMix map[string]int) []string {
	preferred := []string{string(tmux.AgentClaude), string(tmux.AgentCodex), string(tmux.AgentGemini)}
	seen := make(map[string]bool, len(agentMix))
	var order []string
	for _, agent := range preferred {
		if agentMix[agent] > 0 {
			order = append(order, agent)
			seen[agent] = true
		}
	}

	var extras []string
	for agent := range agentMix {
		if seen[agent] {
			continue
		}
		extras = append(extras, agent)
	}
	sort.Strings(extras)
	order = append(order, extras...)
	return order
}

func mergeBudgetDefaults(current, defaults BudgetConfig) BudgetConfig {
	if current.MaxTokensPerMode == 0 {
		current.MaxTokensPerMode = defaults.MaxTokensPerMode
	}
	if current.MaxTotalTokens == 0 {
		current.MaxTotalTokens = defaults.MaxTotalTokens
	}
	if current.TimeoutPerMode == 0 {
		current.TimeoutPerMode = defaults.TimeoutPerMode
	}
	if current.TotalTimeout == 0 {
		current.TotalTimeout = defaults.TotalTimeout
	}
	if current.MaxRetries == 0 {
		current.MaxRetries = defaults.MaxRetries
	}
	return current
}

func mergeSynthesisDefaults(current, defaults SynthesisConfig) SynthesisConfig {
	if current.Strategy == "" {
		current.Strategy = defaults.Strategy
	}
	if current.MinConfidence == 0 {
		current.MinConfidence = defaults.MinConfidence
	}
	if current.MaxFindings == 0 {
		current.MaxFindings = defaults.MaxFindings
	}
	if current.ConflictResolution == "" {
		current.ConflictResolution = defaults.ConflictResolution
	}
	return current
}

func sendPrompt(client *tmux.Client, target, agentType, prompt string, enterDelay, doubleEnterDelay time.Duration) error {
	if err := client.PasteKeys(target, prompt, false); err != nil {
		return fmt.Errorf("send prompt text: %w", err)
	}
	if enterDelay > 0 {
		time.Sleep(enterDelay)
	}
	if err := client.SendKeys(target, "", true); err != nil {
		return fmt.Errorf("send enter: %w", err)
	}
	if needsDoubleEnter(agentType) {
		if doubleEnterDelay > 0 {
			time.Sleep(doubleEnterDelay)
		}
		if err := client.SendKeys(target, "", true); err != nil {
			return fmt.Errorf("send second enter: %w", err)
		}
	}
	return nil
}

func needsDoubleEnter(agentType string) bool {
	switch agentType {
	case "cod", "codex":
		return true
	case "gmi", "gemini":
		return true
	default:
		return false
	}
}
