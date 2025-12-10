// Package cass provides a Go client for CASS (Cross-Agent Session Search).
// CASS indexes prior agent conversations, enabling search across Claude Code,
// Codex, Cursor, Gemini, ChatGPT and other AI coding agent sessions.
package cass

import "time"

// SearchHit represents a single search result from CASS.
type SearchHit struct {
	// SourcePath is the path to the conversation file.
	SourcePath string `json:"source_path"`
	// LineNumber is the line number in the source file (if applicable).
	LineNumber *int `json:"line_number,omitempty"`
	// Agent is the type of AI agent (claude-code, codex, cursor, etc.).
	Agent string `json:"agent"`
	// Workspace is the workspace/project associated with the session.
	Workspace string `json:"workspace"`
	// Title is the conversation title or summary.
	Title string `json:"title"`
	// Score is the relevance score for this result.
	Score float64 `json:"score"`
	// Snippet is a preview of the matching content.
	Snippet string `json:"snippet"`
	// CreatedAt is the Unix timestamp when the session was created.
	CreatedAt *int64 `json:"created_at,omitempty"`
	// MatchType indicates how the result was matched (semantic, keyword, etc.).
	MatchType string `json:"match_type"`
	// Content is the full matching content (if requested).
	Content string `json:"content,omitempty"`
}

// CreatedAtTime returns the CreatedAt field as a time.Time, or zero time if nil.
func (h *SearchHit) CreatedAtTime() time.Time {
	if h.CreatedAt == nil {
		return time.Time{}
	}
	return time.Unix(*h.CreatedAt, 0)
}

// CacheStats contains cache performance statistics.
type CacheStats struct {
	// Hits is the number of cache hits.
	Hits int `json:"hits"`
	// Misses is the number of cache misses.
	Misses int `json:"misses"`
	// Size is the current cache size.
	Size int `json:"size"`
}

// IndexFreshness contains information about index recency.
type IndexFreshness struct {
	// LastUpdate is when the index was last updated.
	LastUpdate string `json:"last_update"`
	// PendingDocs is the number of documents awaiting indexing.
	PendingDocs int `json:"pending_docs"`
	// IsStale indicates if the index is considered stale.
	IsStale bool `json:"is_stale"`
}

// Meta contains response metadata from CASS.
type Meta struct {
	// ElapsedMs is the query execution time in milliseconds.
	ElapsedMs int `json:"elapsed_ms"`
	// WildcardFallback indicates if a wildcard search was used as fallback.
	WildcardFallback bool `json:"wildcard_fallback"`
	// CacheStats contains cache performance information.
	CacheStats *CacheStats `json:"cache_stats,omitempty"`
	// NextCursor is a pagination cursor for fetching more results.
	NextCursor string `json:"next_cursor,omitempty"`
	// RequestID is a unique identifier for this request.
	RequestID string `json:"request_id,omitempty"`
	// IndexFreshness contains index recency information.
	IndexFreshness *IndexFreshness `json:"index_freshness,omitempty"`
	// Warnings contains any warnings from the query.
	Warnings []string `json:"warnings,omitempty"`
}

// HasMore returns true if there are more results available via pagination.
func (m *Meta) HasMore() bool {
	return m != nil && m.NextCursor != ""
}

// AgentCount represents count of results by agent type.
type AgentCount struct {
	Agent string `json:"agent"`
	Count int    `json:"count"`
}

// WorkspaceCount represents count of results by workspace.
type WorkspaceCount struct {
	Workspace string `json:"workspace"`
	Count     int    `json:"count"`
}

// Aggregations contains aggregated statistics about search results.
type Aggregations struct {
	// ByAgent breaks down results by agent type.
	ByAgent []AgentCount `json:"by_agent,omitempty"`
	// ByWorkspace breaks down results by workspace.
	ByWorkspace []WorkspaceCount `json:"by_workspace,omitempty"`
	// TotalSessions is the total number of unique sessions matched.
	TotalSessions int `json:"total_sessions,omitempty"`
	// DateRange contains the date range of matched results.
	DateRange *DateRange `json:"date_range,omitempty"`
}

// DateRange represents a range of dates.
type DateRange struct {
	// Earliest is the earliest date in the range.
	Earliest string `json:"earliest"`
	// Latest is the latest date in the range.
	Latest string `json:"latest"`
}

// SearchResponse represents the response from a CASS search query.
type SearchResponse struct {
	// Query is the original search query.
	Query string `json:"query"`
	// Limit is the maximum number of results requested.
	Limit int `json:"limit"`
	// Offset is the starting offset for pagination.
	Offset int `json:"offset"`
	// Count is the number of results returned in this response.
	Count int `json:"count"`
	// TotalMatches is the total number of matching documents.
	TotalMatches int `json:"total_matches"`
	// Hits contains the search results.
	Hits []SearchHit `json:"hits"`
	// Meta contains response metadata.
	Meta *Meta `json:"_meta,omitempty"`
	// Aggregations contains aggregated statistics.
	Aggregations *Aggregations `json:"aggregations,omitempty"`
}

// HasResults returns true if there are any search results.
func (r *SearchResponse) HasResults() bool {
	return len(r.Hits) > 0
}

// HasMore returns true if there are more results available.
func (r *SearchResponse) HasMore() bool {
	return r.Meta.HasMore() || r.Offset+r.Count < r.TotalMatches
}

// IndexInfo contains information about the CASS search index.
type IndexInfo struct {
	// DocCount is the total number of indexed documents.
	DocCount int `json:"doc_count"`
	// SizeBytes is the index size in bytes.
	SizeBytes int64 `json:"size_bytes"`
	// LastUpdated is when the index was last updated.
	LastUpdated string `json:"last_updated"`
	// Healthy indicates if the index is healthy.
	Healthy bool `json:"healthy"`
}

// SizeMB returns the index size in megabytes.
func (i *IndexInfo) SizeMB() float64 {
	return float64(i.SizeBytes) / (1024 * 1024)
}

// DBInfo contains information about the CASS database.
type DBInfo struct {
	// Path is the database file path.
	Path string `json:"path"`
	// SizeBytes is the database size in bytes.
	SizeBytes int64 `json:"size_bytes"`
	// Healthy indicates if the database is healthy.
	Healthy bool `json:"healthy"`
	// SessionCount is the number of sessions in the database.
	SessionCount int `json:"session_count"`
}

// SizeMB returns the database size in megabytes.
func (d *DBInfo) SizeMB() float64 {
	return float64(d.SizeBytes) / (1024 * 1024)
}

// Pending contains information about pending operations.
type Pending struct {
	// Sessions is the number of sessions awaiting processing.
	Sessions int `json:"sessions"`
	// Files is the number of files awaiting indexing.
	Files int `json:"files"`
}

// HasPending returns true if there are any pending operations.
func (p *Pending) HasPending() bool {
	return p.Sessions > 0 || p.Files > 0
}

// StatusResponse represents the CASS health status response.
type StatusResponse struct {
	// Healthy indicates overall system health.
	Healthy bool `json:"healthy"`
	// RecommendedAction is a suggested action if unhealthy.
	RecommendedAction string `json:"recommended_action"`
	// Index contains index status information.
	Index IndexInfo `json:"index"`
	// Database contains database status information.
	Database DBInfo `json:"database"`
	// Pending contains pending operations information.
	Pending Pending `json:"pending"`
	// Meta contains response metadata.
	Meta *Meta `json:"_meta,omitempty"`
}

// IsHealthy returns true if CASS is fully healthy.
func (s *StatusResponse) IsHealthy() bool {
	return s.Healthy && s.Index.Healthy && s.Database.Healthy
}

// NeedsAction returns true if a recommended action exists.
func (s *StatusResponse) NeedsAction() bool {
	return s.RecommendedAction != ""
}

// Limits contains CASS operational limits.
type Limits struct {
	// MaxQueryLength is the maximum query string length.
	MaxQueryLength int `json:"max_query_length"`
	// MaxResults is the maximum number of results per query.
	MaxResults int `json:"max_results"`
	// MaxConcurrentQueries is the maximum concurrent queries.
	MaxConcurrentQueries int `json:"max_concurrent_queries"`
	// RateLimitPerMinute is the rate limit per minute.
	RateLimitPerMinute int `json:"rate_limit_per_minute"`
}

// Capabilities represents CASS feature discovery response.
type Capabilities struct {
	// CrateVersion is the CASS crate version.
	CrateVersion string `json:"crate_version"`
	// APIVersion is the API version number.
	APIVersion int `json:"api_version"`
	// ContractVersion is the contract/schema version.
	ContractVersion string `json:"contract_version"`
	// Features is the list of enabled features.
	Features []string `json:"features"`
	// Connectors is the list of available connectors (claude-code, codex, etc.).
	Connectors []string `json:"connectors"`
	// Limits contains operational limits.
	Limits Limits `json:"limits"`
}

// HasFeature checks if a specific feature is enabled.
func (c *Capabilities) HasFeature(feature string) bool {
	for _, f := range c.Features {
		if f == feature {
			return true
		}
	}
	return false
}

// HasConnector checks if a specific connector is available.
func (c *Capabilities) HasConnector(connector string) bool {
	for _, conn := range c.Connectors {
		if conn == connector {
			return true
		}
	}
	return false
}

// CASSError represents an error from CASS.
type CASSError struct {
	// Code is the error code.
	Code int `json:"code"`
	// Kind is the error type/category.
	Kind string `json:"kind"`
	// Message is the human-readable error message.
	Message string `json:"message"`
	// Hint is an optional hint for resolving the error.
	Hint string `json:"hint,omitempty"`
	// Retryable indicates if the operation can be retried.
	Retryable bool `json:"retryable"`
}

// Error implements the error interface.
func (e *CASSError) Error() string {
	if e.Hint != "" {
		return e.Message + " (hint: " + e.Hint + ")"
	}
	return e.Message
}

// ErrorResponse wraps a CASS error response.
type ErrorResponse struct {
	Error CASSError `json:"error"`
}

// ViewResponse represents the response from viewing a specific conversation.
type ViewResponse struct {
	// SourcePath is the path to the conversation file.
	SourcePath string `json:"source_path"`
	// Agent is the type of AI agent.
	Agent string `json:"agent"`
	// Workspace is the workspace/project.
	Workspace string `json:"workspace"`
	// Title is the conversation title.
	Title string `json:"title"`
	// CreatedAt is when the session was created.
	CreatedAt *int64 `json:"created_at,omitempty"`
	// Messages contains the conversation messages.
	Messages []Message `json:"messages"`
	// Meta contains response metadata.
	Meta *Meta `json:"_meta,omitempty"`
}

// Message represents a single message in a conversation.
type Message struct {
	// Role is the message role (user, assistant, system).
	Role string `json:"role"`
	// Content is the message content.
	Content string `json:"content"`
	// Timestamp is when the message was sent.
	Timestamp *int64 `json:"timestamp,omitempty"`
	// LineNumber is the line number in the source file.
	LineNumber int `json:"line_number,omitempty"`
}

// TimestampTime returns the Timestamp field as a time.Time, or zero time if nil.
func (m *Message) TimestampTime() time.Time {
	if m.Timestamp == nil {
		return time.Time{}
	}
	return time.Unix(*m.Timestamp, 0)
}

// ExpandResponse represents an expanded context view around a match.
type ExpandResponse struct {
	// SourcePath is the path to the conversation file.
	SourcePath string `json:"source_path"`
	// CenterLine is the line number that was expanded around.
	CenterLine int `json:"center_line"`
	// ContextBefore is the number of lines before.
	ContextBefore int `json:"context_before"`
	// ContextAfter is the number of lines after.
	ContextAfter int `json:"context_after"`
	// Messages contains the expanded context messages.
	Messages []Message `json:"messages"`
	// Meta contains response metadata.
	Meta *Meta `json:"_meta,omitempty"`
}

// TimelineEntry represents an entry in the activity timeline.
type TimelineEntry struct {
	// Timestamp is when the activity occurred.
	Timestamp int64 `json:"timestamp"`
	// Agent is the type of AI agent.
	Agent string `json:"agent"`
	// Workspace is the workspace/project.
	Workspace string `json:"workspace"`
	// Title is the session title.
	Title string `json:"title"`
	// SourcePath is the path to the session file.
	SourcePath string `json:"source_path"`
	// MessageCount is the number of messages in the session.
	MessageCount int `json:"message_count"`
}

// TimestampTime returns the Timestamp as a time.Time.
func (e *TimelineEntry) TimestampTime() time.Time {
	return time.Unix(e.Timestamp, 0)
}

// TimelineResponse represents the response from the timeline API.
type TimelineResponse struct {
	// Entries contains the timeline entries.
	Entries []TimelineEntry `json:"entries"`
	// Count is the number of entries returned.
	Count int `json:"count"`
	// Meta contains response metadata.
	Meta *Meta `json:"_meta,omitempty"`
}
