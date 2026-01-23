package ensemble

import (
	"testing"
	"time"
)

func TestModeCategory_IsValid(t *testing.T) {
	tests := []struct {
		cat   ModeCategory
		valid bool
	}{
		{CategoryFormal, true},
		{CategoryAmpliative, true},
		{CategoryUncertainty, true},
		{CategoryVagueness, true},
		{CategoryChange, true},
		{CategoryCausal, true},
		{CategoryPractical, true},
		{CategoryStrategic, true},
		{CategoryDialectical, true},
		{CategoryModal, true},
		{CategoryDomain, true},
		{CategoryMeta, true},
		{ModeCategory("invalid"), false},
		{ModeCategory(""), false},
	}

	for _, tt := range tests {
		if got := tt.cat.IsValid(); got != tt.valid {
			t.Errorf("ModeCategory(%q).IsValid() = %v, want %v", tt.cat, got, tt.valid)
		}
	}
}

func TestValidateModeID(t *testing.T) {
	tests := []struct {
		id      string
		wantErr bool
	}{
		{"deductive", false},
		{"bayesian-inference", false},
		{"foo-bar-baz", false},
		{"a", false},
		{"a1", false},
		{"", true},                       // empty
		{"Deductive", true},              // uppercase
		{"123abc", true},                 // starts with number
		{"-invalid", true},               // starts with hyphen
		{"has spaces", true},             // contains spaces
		{"has_underscore", true},         // contains underscore
		{string(make([]byte, 65)), true}, // too long
	}

	for _, tt := range tests {
		err := ValidateModeID(tt.id)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateModeID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
		}
	}
}

func TestReasoningMode_Validate(t *testing.T) {
	validMode := ReasoningMode{
		ID:        "deductive",
		Name:      "Deductive Logic",
		Category:  CategoryFormal,
		ShortDesc: "Derive conclusions from premises",
	}

	if err := validMode.Validate(); err != nil {
		t.Errorf("valid mode should pass validation: %v", err)
	}

	// Test missing ID
	noID := validMode
	noID.ID = ""
	if err := noID.Validate(); err == nil {
		t.Error("mode without ID should fail validation")
	}

	// Test invalid ID
	invalidID := validMode
	invalidID.ID = "INVALID"
	if err := invalidID.Validate(); err == nil {
		t.Error("mode with invalid ID should fail validation")
	}

	// Test missing name
	noName := validMode
	noName.Name = ""
	if err := noName.Validate(); err == nil {
		t.Error("mode without name should fail validation")
	}

	// Test invalid category
	invalidCat := validMode
	invalidCat.Category = "invalid"
	if err := invalidCat.Validate(); err == nil {
		t.Error("mode with invalid category should fail validation")
	}

	// Test long short_desc
	longDesc := validMode
	longDesc.ShortDesc = string(make([]byte, 100))
	if err := longDesc.Validate(); err == nil {
		t.Error("mode with short_desc > 80 chars should fail validation")
	}
}

func TestAssignmentStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status   AssignmentStatus
		terminal bool
	}{
		{AssignmentPending, false},
		{AssignmentInjecting, false},
		{AssignmentActive, false},
		{AssignmentDone, true},
		{AssignmentError, true},
	}

	for _, tt := range tests {
		if got := tt.status.IsTerminal(); got != tt.terminal {
			t.Errorf("AssignmentStatus(%q).IsTerminal() = %v, want %v", tt.status, got, tt.terminal)
		}
	}
}

func TestEnsembleStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status   EnsembleStatus
		terminal bool
	}{
		{EnsembleSpawning, false},
		{EnsembleInjecting, false},
		{EnsembleActive, false},
		{EnsembleSynthesizing, false},
		{EnsembleComplete, true},
		{EnsembleError, true},
	}

	for _, tt := range tests {
		if got := tt.status.IsTerminal(); got != tt.terminal {
			t.Errorf("EnsembleStatus(%q).IsTerminal() = %v, want %v", tt.status, got, tt.terminal)
		}
	}
}

func TestSynthesisStrategy_IsValid(t *testing.T) {
	tests := []struct {
		strategy SynthesisStrategy
		valid    bool
	}{
		{StrategyConsensus, true},
		{StrategyDebate, true},
		{StrategyWeighted, true},
		{StrategySequential, true},
		{StrategyBestOf, true},
		{SynthesisStrategy("invalid"), false},
		{SynthesisStrategy(""), false},
	}

	for _, tt := range tests {
		if got := tt.strategy.IsValid(); got != tt.valid {
			t.Errorf("SynthesisStrategy(%q).IsValid() = %v, want %v", tt.strategy, got, tt.valid)
		}
	}
}

func TestEnsemblePreset_Validate(t *testing.T) {
	catalog := []ReasoningMode{
		{ID: "deductive", Name: "Deductive", Category: CategoryFormal, ShortDesc: "Test"},
		{ID: "bayesian", Name: "Bayesian", Category: CategoryUncertainty, ShortDesc: "Test"},
	}

	validPreset := EnsemblePreset{
		Name:              "test-preset",
		Description:       "A test preset",
		Modes:             []string{"deductive", "bayesian"},
		SynthesisStrategy: StrategyConsensus,
	}

	if err := validPreset.Validate(catalog); err != nil {
		t.Errorf("valid preset should pass validation: %v", err)
	}

	// Test missing name
	noName := validPreset
	noName.Name = ""
	if err := noName.Validate(catalog); err == nil {
		t.Error("preset without name should fail validation")
	}

	// Test no modes
	noModes := validPreset
	noModes.Modes = nil
	if err := noModes.Validate(catalog); err == nil {
		t.Error("preset without modes should fail validation")
	}

	// Test invalid strategy
	invalidStrategy := validPreset
	invalidStrategy.SynthesisStrategy = "invalid"
	if err := invalidStrategy.Validate(catalog); err == nil {
		t.Error("preset with invalid strategy should fail validation")
	}

	// Test missing mode
	missingMode := validPreset
	missingMode.Modes = []string{"deductive", "nonexistent"}
	if err := missingMode.Validate(catalog); err == nil {
		t.Error("preset referencing nonexistent mode should fail validation")
	}
}

func TestModeCatalog(t *testing.T) {
	modes := []ReasoningMode{
		{ID: "deductive", Name: "Deductive Logic", Category: CategoryFormal, ShortDesc: "Derive conclusions", BestFor: []string{"proofs"}},
		{ID: "bayesian", Name: "Bayesian Inference", Category: CategoryUncertainty, ShortDesc: "Probabilistic reasoning", BestFor: []string{"prediction"}},
		{ID: "causal-inference", Name: "Causal Inference", Category: CategoryCausal, ShortDesc: "Find causes", BestFor: []string{"debugging"}},
	}

	cat, err := NewModeCatalog(modes, "1.0.0")
	if err != nil {
		t.Fatalf("NewModeCatalog failed: %v", err)
	}

	// Test Count
	if got := cat.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}

	// Test Version
	if got := cat.Version(); got != "1.0.0" {
		t.Errorf("Version() = %q, want %q", got, "1.0.0")
	}

	// Test GetMode
	mode := cat.GetMode("deductive")
	if mode == nil {
		t.Error("GetMode(deductive) returned nil")
	} else if mode.Name != "Deductive Logic" {
		t.Errorf("GetMode(deductive).Name = %q, want %q", mode.Name, "Deductive Logic")
	}

	// Test GetMode nonexistent
	if cat.GetMode("nonexistent") != nil {
		t.Error("GetMode(nonexistent) should return nil")
	}

	// Test ListModes
	all := cat.ListModes()
	if len(all) != 3 {
		t.Errorf("ListModes() returned %d modes, want 3", len(all))
	}

	// Test ListByCategory
	formal := cat.ListByCategory(CategoryFormal)
	if len(formal) != 1 {
		t.Errorf("ListByCategory(Formal) returned %d modes, want 1", len(formal))
	}

	// Test SearchModes
	found := cat.SearchModes("logic")
	if len(found) != 1 {
		t.Errorf("SearchModes(logic) returned %d modes, want 1", len(found))
	}

	// Test search in BestFor
	foundBestFor := cat.SearchModes("proofs")
	if len(foundBestFor) != 1 {
		t.Errorf("SearchModes(proofs) returned %d modes, want 1", len(foundBestFor))
	}
}

func TestModeCatalog_DuplicateID(t *testing.T) {
	modes := []ReasoningMode{
		{ID: "deductive", Name: "Deductive Logic", Category: CategoryFormal, ShortDesc: "Test"},
		{ID: "deductive", Name: "Duplicate", Category: CategoryFormal, ShortDesc: "Test"},
	}

	_, err := NewModeCatalog(modes, "1.0.0")
	if err == nil {
		t.Error("NewModeCatalog should fail with duplicate IDs")
	}
}

func TestModeCatalog_InvalidMode(t *testing.T) {
	modes := []ReasoningMode{
		{ID: "INVALID", Name: "Invalid", Category: CategoryFormal, ShortDesc: "Test"},
	}

	_, err := NewModeCatalog(modes, "1.0.0")
	if err == nil {
		t.Error("NewModeCatalog should fail with invalid mode")
	}
}

func TestModeAssignment_Fields(t *testing.T) {
	now := time.Now()
	completed := now.Add(time.Hour)

	assignment := ModeAssignment{
		ModeID:      "deductive",
		PaneName:    "myproject__cc_1",
		AgentType:   "cc",
		Status:      AssignmentActive,
		OutputPath:  "/tmp/output.txt",
		AssignedAt:  now,
		CompletedAt: &completed,
		Error:       "",
	}

	if assignment.ModeID != "deductive" {
		t.Errorf("ModeID = %q, want %q", assignment.ModeID, "deductive")
	}
	if assignment.Status != AssignmentActive {
		t.Errorf("Status = %q, want %q", assignment.Status, AssignmentActive)
	}
}

func TestEnsembleSession_Fields(t *testing.T) {
	now := time.Now()

	session := EnsembleSession{
		SessionName:       "myproject",
		Question:          "What is the best approach?",
		PresetUsed:        "architecture-review",
		Assignments:       []ModeAssignment{},
		Status:            EnsembleActive,
		SynthesisStrategy: StrategyConsensus,
		CreatedAt:         now,
	}

	if session.Status != EnsembleActive {
		t.Errorf("Status = %q, want %q", session.Status, EnsembleActive)
	}
	if session.SynthesisStrategy != StrategyConsensus {
		t.Errorf("SynthesisStrategy = %q, want %q", session.SynthesisStrategy, StrategyConsensus)
	}
}

func TestAllCategories(t *testing.T) {
	cats := AllCategories()
	if len(cats) != 12 {
		t.Errorf("AllCategories() returned %d categories, want 12", len(cats))
	}

	// All should be valid
	for _, cat := range cats {
		if !cat.IsValid() {
			t.Errorf("AllCategories() returned invalid category %q", cat)
		}
	}
}

// =============================================================================
// Output Schema Tests
// =============================================================================

func TestImpactLevel_IsValid(t *testing.T) {
	tests := []struct {
		level ImpactLevel
		valid bool
	}{
		{ImpactHigh, true},
		{ImpactMedium, true},
		{ImpactLow, true},
		{ImpactLevel("invalid"), false},
		{ImpactLevel(""), false},
	}

	for _, tt := range tests {
		if got := tt.level.IsValid(); got != tt.valid {
			t.Errorf("ImpactLevel(%q).IsValid() = %v, want %v", tt.level, got, tt.valid)
		}
	}
}

func TestConfidence_Validate(t *testing.T) {
	tests := []struct {
		conf    Confidence
		wantErr bool
	}{
		{0.0, false},
		{0.5, false},
		{1.0, false},
		{0.75, false},
		{-0.1, true},
		{1.1, true},
		{2.0, true},
	}

	for _, tt := range tests {
		err := tt.conf.Validate()
		if (err != nil) != tt.wantErr {
			t.Errorf("Confidence(%v).Validate() error = %v, wantErr %v", tt.conf, err, tt.wantErr)
		}
	}
}

func TestConfidence_String(t *testing.T) {
	tests := []struct {
		conf Confidence
		want string
	}{
		{0.0, "0%"},
		{0.5, "50%"},
		{1.0, "100%"},
		{0.75, "75%"},
	}

	for _, tt := range tests {
		if got := tt.conf.String(); got != tt.want {
			t.Errorf("Confidence(%v).String() = %q, want %q", tt.conf, got, tt.want)
		}
	}
}

func TestFinding_Validate(t *testing.T) {
	validFinding := Finding{
		Finding:    "Test finding",
		Impact:     ImpactHigh,
		Confidence: 0.8,
	}

	if err := validFinding.Validate(); err != nil {
		t.Errorf("valid finding should pass validation: %v", err)
	}

	// Test missing finding text
	noText := validFinding
	noText.Finding = ""
	if err := noText.Validate(); err == nil {
		t.Error("finding without text should fail validation")
	}

	// Test invalid impact
	invalidImpact := validFinding
	invalidImpact.Impact = "invalid"
	if err := invalidImpact.Validate(); err == nil {
		t.Error("finding with invalid impact should fail validation")
	}

	// Test invalid confidence
	invalidConf := validFinding
	invalidConf.Confidence = 1.5
	if err := invalidConf.Validate(); err == nil {
		t.Error("finding with invalid confidence should fail validation")
	}
}

func TestRisk_Validate(t *testing.T) {
	validRisk := Risk{
		Risk:       "Test risk",
		Impact:     ImpactMedium,
		Likelihood: 0.3,
	}

	if err := validRisk.Validate(); err != nil {
		t.Errorf("valid risk should pass validation: %v", err)
	}

	// Test missing risk text
	noText := validRisk
	noText.Risk = ""
	if err := noText.Validate(); err == nil {
		t.Error("risk without text should fail validation")
	}

	// Test invalid impact
	invalidImpact := validRisk
	invalidImpact.Impact = "invalid"
	if err := invalidImpact.Validate(); err == nil {
		t.Error("risk with invalid impact should fail validation")
	}

	// Test invalid likelihood
	invalidLikelihood := validRisk
	invalidLikelihood.Likelihood = -0.5
	if err := invalidLikelihood.Validate(); err == nil {
		t.Error("risk with invalid likelihood should fail validation")
	}
}

func TestRecommendation_Validate(t *testing.T) {
	validRec := Recommendation{
		Recommendation: "Test recommendation",
		Priority:       ImpactHigh,
	}

	if err := validRec.Validate(); err != nil {
		t.Errorf("valid recommendation should pass validation: %v", err)
	}

	// Test missing text
	noText := validRec
	noText.Recommendation = ""
	if err := noText.Validate(); err == nil {
		t.Error("recommendation without text should fail validation")
	}

	// Test invalid priority
	invalidPriority := validRec
	invalidPriority.Priority = "invalid"
	if err := invalidPriority.Validate(); err == nil {
		t.Error("recommendation with invalid priority should fail validation")
	}
}

func TestQuestion_Validate(t *testing.T) {
	validQuestion := Question{
		Question: "What is the requirement?",
	}

	if err := validQuestion.Validate(); err != nil {
		t.Errorf("valid question should pass validation: %v", err)
	}

	// Test missing question text
	noText := validQuestion
	noText.Question = ""
	if err := noText.Validate(); err == nil {
		t.Error("question without text should fail validation")
	}
}

func TestFailureModeWarning_Validate(t *testing.T) {
	validWarning := FailureModeWarning{
		Mode:        "confirmation-bias",
		Description: "Seeking evidence that confirms existing beliefs",
	}

	if err := validWarning.Validate(); err != nil {
		t.Errorf("valid failure mode warning should pass validation: %v", err)
	}

	// Test missing mode
	noMode := validWarning
	noMode.Mode = ""
	if err := noMode.Validate(); err == nil {
		t.Error("warning without mode should fail validation")
	}

	// Test missing description
	noDesc := validWarning
	noDesc.Description = ""
	if err := noDesc.Validate(); err == nil {
		t.Error("warning without description should fail validation")
	}
}

func TestModeOutput_Validate(t *testing.T) {
	validOutput := ModeOutput{
		ModeID: "deductive",
		Thesis: "The system has a critical bug in the auth module",
		TopFindings: []Finding{
			{
				Finding:    "Missing input validation",
				Impact:     ImpactHigh,
				Confidence: 0.9,
			},
		},
		Confidence:  0.85,
		GeneratedAt: time.Now(),
	}

	if err := validOutput.Validate(); err != nil {
		t.Errorf("valid mode output should pass validation: %v", err)
	}

	// Test missing mode_id
	noModeID := validOutput
	noModeID.ModeID = ""
	if err := noModeID.Validate(); err == nil {
		t.Error("output without mode_id should fail validation")
	}

	// Test missing thesis
	noThesis := validOutput
	noThesis.Thesis = ""
	if err := noThesis.Validate(); err == nil {
		t.Error("output without thesis should fail validation")
	}

	// Test no findings
	noFindings := validOutput
	noFindings.TopFindings = nil
	if err := noFindings.Validate(); err == nil {
		t.Error("output without findings should fail validation")
	}

	// Test invalid confidence
	invalidConf := validOutput
	invalidConf.Confidence = 1.5
	if err := invalidConf.Validate(); err == nil {
		t.Error("output with invalid confidence should fail validation")
	}

	// Test with invalid finding
	invalidFinding := validOutput
	invalidFinding.TopFindings = []Finding{{Finding: "", Impact: ImpactHigh, Confidence: 0.5}}
	if err := invalidFinding.Validate(); err == nil {
		t.Error("output with invalid finding should fail validation")
	}
}

func TestModeOutput_ValidateNestedTypes(t *testing.T) {
	validOutput := ModeOutput{
		ModeID:     "test",
		Thesis:     "Test thesis",
		Confidence: 0.5,
		TopFindings: []Finding{
			{Finding: "Valid finding", Impact: ImpactHigh, Confidence: 0.8},
		},
		Risks: []Risk{
			{Risk: "Valid risk", Impact: ImpactMedium, Likelihood: 0.3},
		},
		Recommendations: []Recommendation{
			{Recommendation: "Valid rec", Priority: ImpactLow},
		},
		QuestionsForUser: []Question{
			{Question: "Valid question?"},
		},
		FailureModesToWatch: []FailureModeWarning{
			{Mode: "bias", Description: "Confirmation bias"},
		},
		GeneratedAt: time.Now(),
	}

	if err := validOutput.Validate(); err != nil {
		t.Errorf("output with all valid nested types should pass: %v", err)
	}

	// Test invalid risk
	invalidRisk := validOutput
	invalidRisk.Risks = []Risk{{Risk: "", Impact: ImpactHigh, Likelihood: 0.5}}
	if err := invalidRisk.Validate(); err == nil {
		t.Error("output with invalid risk should fail validation")
	}

	// Test invalid recommendation
	invalidRec := validOutput
	invalidRec.Recommendations = []Recommendation{{Recommendation: "", Priority: ImpactHigh}}
	if err := invalidRec.Validate(); err == nil {
		t.Error("output with invalid recommendation should fail validation")
	}

	// Test invalid question
	invalidQ := validOutput
	invalidQ.QuestionsForUser = []Question{{Question: ""}}
	if err := invalidQ.Validate(); err == nil {
		t.Error("output with invalid question should fail validation")
	}

	// Test invalid failure mode warning
	invalidFM := validOutput
	invalidFM.FailureModesToWatch = []FailureModeWarning{{Mode: "", Description: ""}}
	if err := invalidFM.Validate(); err == nil {
		t.Error("output with invalid failure mode should fail validation")
	}
}

func TestDefaultBudgetConfig(t *testing.T) {
	cfg := DefaultBudgetConfig()

	if cfg.MaxTokensPerMode <= 0 {
		t.Error("MaxTokensPerMode should be positive")
	}
	if cfg.MaxTotalTokens <= 0 {
		t.Error("MaxTotalTokens should be positive")
	}
	if cfg.TimeoutPerMode <= 0 {
		t.Error("TimeoutPerMode should be positive")
	}
	if cfg.TotalTimeout <= 0 {
		t.Error("TotalTimeout should be positive")
	}
}

func TestDefaultSynthesisConfig(t *testing.T) {
	cfg := DefaultSynthesisConfig()

	if !cfg.Strategy.IsValid() {
		t.Error("Strategy should be valid")
	}
	if cfg.MinConfidence < 0 || cfg.MinConfidence > 1 {
		t.Errorf("MinConfidence should be between 0 and 1, got %v", cfg.MinConfidence)
	}
	if cfg.MaxFindings <= 0 {
		t.Error("MaxFindings should be positive")
	}
}

func TestEnsemble_Validate(t *testing.T) {
	modes := []ReasoningMode{
		{ID: "deductive", Name: "Deductive", Category: CategoryFormal, ShortDesc: "Test"},
		{ID: "bayesian", Name: "Bayesian", Category: CategoryUncertainty, ShortDesc: "Test"},
	}
	catalog, _ := NewModeCatalog(modes, "1.0.0")

	validEnsemble := Ensemble{
		Name:        "test-ensemble",
		DisplayName: "Test Ensemble",
		Description: "A test ensemble",
		ModeIDs:     []string{"deductive", "bayesian"},
		Synthesis:   DefaultSynthesisConfig(),
		Budget:      DefaultBudgetConfig(),
	}

	if err := validEnsemble.Validate(catalog); err != nil {
		t.Errorf("valid ensemble should pass validation: %v", err)
	}

	// Test missing name
	noName := validEnsemble
	noName.Name = ""
	if err := noName.Validate(catalog); err == nil {
		t.Error("ensemble without name should fail validation")
	}

	// Test invalid name format
	invalidName := validEnsemble
	invalidName.Name = "INVALID"
	if err := invalidName.Validate(catalog); err == nil {
		t.Error("ensemble with invalid name format should fail validation")
	}

	// Test missing display name
	noDisplayName := validEnsemble
	noDisplayName.DisplayName = ""
	if err := noDisplayName.Validate(catalog); err == nil {
		t.Error("ensemble without display_name should fail validation")
	}

	// Test no modes
	noModes := validEnsemble
	noModes.ModeIDs = nil
	if err := noModes.Validate(catalog); err == nil {
		t.Error("ensemble without modes should fail validation")
	}

	// Test nonexistent mode
	badMode := validEnsemble
	badMode.ModeIDs = []string{"deductive", "nonexistent"}
	if err := badMode.Validate(catalog); err == nil {
		t.Error("ensemble with nonexistent mode should fail validation")
	}
}
