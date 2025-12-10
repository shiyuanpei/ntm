package cass

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// Common errors
var (
	// ErrNotInstalled indicates CASS is not installed or not in PATH.
	ErrNotInstalled = errors.New("cass is not installed or not in PATH")
	// ErrTimeout indicates the command timed out.
	ErrTimeout = errors.New("cass command timed out")
	// ErrNoResults indicates no results were found.
	ErrNoResults = errors.New("no results found")
)

const (
	// DefaultTimeout is the default command timeout.
	DefaultTimeout = 30 * time.Second
	// DefaultSearchLimit is the default number of search results.
	DefaultSearchLimit = 10
)

// Client provides methods to interact with CASS via its robot mode CLI.
type Client struct {
	// Timeout is the command execution timeout.
	Timeout time.Duration
	// cassPath is the path to the cass binary (empty = use PATH).
	cassPath string
}

// Option configures the Client.
type Option func(*Client)

// WithTimeout sets the command timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// WithCASSPath sets an explicit path to the cass binary.
func WithCASSPath(path string) Option {
	return func(c *Client) {
		c.cassPath = path
	}
}

// NewClient creates a new CASS client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		Timeout: DefaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// IsInstalled checks if CASS is available in PATH.
func (c *Client) IsInstalled() bool {
	path := c.cassPath
	if path == "" {
		path = "cass"
	}
	_, err := exec.LookPath(path)
	return err == nil
}

// exec runs a CASS command and returns the output.
func (c *Client) exec(ctx context.Context, args ...string) ([]byte, error) {
	path := c.cassPath
	if path == "" {
		path = "cass"
	}

	// Create command with context for timeout
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ErrTimeout
		}

		// Check if cass is not installed
		if errors.Is(err, exec.ErrNotFound) {
			return nil, ErrNotInstalled
		}

		// Try to parse error response from stderr
		if stderr.Len() > 0 {
			var errResp ErrorResponse
			if jsonErr := json.Unmarshal(stderr.Bytes(), &errResp); jsonErr == nil {
				return nil, &errResp.Error
			}
		}

		// Return generic error with stderr content
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("cass error: %s", stderr.String())
		}
		return nil, fmt.Errorf("cass error: %w", err)
	}

	return stdout.Bytes(), nil
}

// Health checks CASS health status.
func (c *Client) Health(ctx context.Context) (*StatusResponse, error) {
	output, err := c.exec(ctx, "health", "--json")
	if err != nil {
		return nil, err
	}

	var status StatusResponse
	if err := json.Unmarshal(output, &status); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &status, nil
}

// Capabilities returns CASS feature discovery information.
func (c *Client) Capabilities(ctx context.Context) (*Capabilities, error) {
	output, err := c.exec(ctx, "capabilities", "--json")
	if err != nil {
		return nil, err
	}

	var caps Capabilities
	if err := json.Unmarshal(output, &caps); err != nil {
		return nil, fmt.Errorf("failed to parse capabilities response: %w", err)
	}

	return &caps, nil
}

// SearchOptions configures a search query.
type SearchOptions struct {
	// Query is the search query string.
	Query string
	// Limit is the maximum number of results (default: 10).
	Limit int
	// Agent filters results by agent type (claude-code, codex, etc.).
	Agent string
	// Workspace filters results by workspace/project.
	Workspace string
	// Days limits search to the last N days.
	Days int
	// Fields specifies which fields to include (minimal, full, etc.).
	Fields string
}

// Search performs a search query against CASS.
func (c *Client) Search(ctx context.Context, opts SearchOptions) (*SearchResponse, error) {
	if opts.Query == "" {
		return nil, errors.New("query is required")
	}

	args := []string{"search", opts.Query, "--robot"}

	if opts.Limit > 0 {
		args = append(args, "--limit", fmt.Sprintf("%d", opts.Limit))
	} else {
		args = append(args, "--limit", fmt.Sprintf("%d", DefaultSearchLimit))
	}

	if opts.Agent != "" {
		args = append(args, "--agent", opts.Agent)
	}
	if opts.Workspace != "" {
		args = append(args, "--workspace", opts.Workspace)
	}
	if opts.Days > 0 {
		args = append(args, "--days", fmt.Sprintf("%d", opts.Days))
	}
	if opts.Fields != "" {
		args = append(args, "--fields", opts.Fields)
	}

	output, err := c.exec(ctx, args...)
	if err != nil {
		return nil, err
	}

	var resp SearchResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	return &resp, nil
}

// View retrieves a specific conversation by source path and optional line number.
func (c *Client) View(ctx context.Context, sourcePath string, lineNumber int) (*ViewResponse, error) {
	if sourcePath == "" {
		return nil, errors.New("source path is required")
	}

	args := []string{"view", sourcePath, "--json"}
	if lineNumber > 0 {
		args = append(args, "-n", fmt.Sprintf("%d", lineNumber))
	}

	output, err := c.exec(ctx, args...)
	if err != nil {
		return nil, err
	}

	var resp ViewResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse view response: %w", err)
	}

	return &resp, nil
}

// Expand retrieves expanded context around a specific line in a conversation.
func (c *Client) Expand(ctx context.Context, sourcePath string, lineNumber, contextLines int) (*ExpandResponse, error) {
	if sourcePath == "" {
		return nil, errors.New("source path is required")
	}
	if lineNumber <= 0 {
		return nil, errors.New("line number is required")
	}

	args := []string{"expand", sourcePath, "-n", fmt.Sprintf("%d", lineNumber), "--json"}
	if contextLines > 0 {
		args = append(args, "-C", fmt.Sprintf("%d", contextLines))
	}

	output, err := c.exec(ctx, args...)
	if err != nil {
		return nil, err
	}

	var resp ExpandResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse expand response: %w", err)
	}

	return &resp, nil
}

// TimelineOptions configures a timeline query.
type TimelineOptions struct {
	// Limit is the maximum number of entries.
	Limit int
	// Agent filters by agent type.
	Agent string
	// Workspace filters by workspace.
	Workspace string
	// Days limits to the last N days.
	Days int
}

// Timeline retrieves the activity timeline.
func (c *Client) Timeline(ctx context.Context, opts TimelineOptions) (*TimelineResponse, error) {
	args := []string{"timeline", "--json"}

	if opts.Limit > 0 {
		args = append(args, "--limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Agent != "" {
		args = append(args, "--agent", opts.Agent)
	}
	if opts.Workspace != "" {
		args = append(args, "--workspace", opts.Workspace)
	}
	if opts.Days > 0 {
		args = append(args, "--days", fmt.Sprintf("%d", opts.Days))
	}

	output, err := c.exec(ctx, args...)
	if err != nil {
		return nil, err
	}

	var resp TimelineResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse timeline response: %w", err)
	}

	return &resp, nil
}
