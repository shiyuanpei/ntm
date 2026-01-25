# NTM Robot Mode - AI Automation Guide

Machine-readable JSON API for AI agents and automation systems.

## Quick Start

```bash
# 自描述 API - 获取所有可用命令
ntm --robot-capabilities | jq '.commands[] | {name, flag, description}'
```

## Core Concepts

### Response Envelope

所有 robot 命令返回统一的 JSON 结构：

```json
{
  "success": true,
  "timestamp": "2026-01-23T07:24:05Z",
  "error": null,
  "error_code": null,
  "hint": null,
  "_agent_hints": {
    "summary": "Human-readable summary",
    "suggested_actions": [...],
    "notes": [...]
  }
}
```

### Agent Types

| Type | Aliases | Description |
|------|---------|-------------|
| claude | cc | Claude Code |
| codex | cod | Codex CLI |
| gemini | gmi | Gemini CLI |
| cursor | - | Cursor |
| windsurf | - | Windsurf |
| aider | - | Aider |

---

## Bead Management (Task Tracking)

### Create Task

```bash
ntm --robot-bead-create \
  --bead-title="Fix login validation" \
  --bead-type=bug \
  --bead-priority=1 \
  --bead-labels="auth,urgent" \
  --bead-description="Login fails on mobile devices"
```

**Response:**
```json
{
  "success": true,
  "bead_id": "bd-abc123",
  "title": "Fix login validation",
  "type": "bug",
  "priority": "P1",
  "created": true
}
```

**Parameters:**
| Param | Required | Default | Values |
|-------|----------|---------|--------|
| `--bead-title` | Yes | - | string |
| `--bead-type` | No | task | task, bug, feature, epic, chore |
| `--bead-priority` | No | 2 | 0-4 (0=critical, 4=backlog) |
| `--bead-labels` | No | - | comma-separated |
| `--bead-description` | No | - | string |
| `--bead-depends-on` | No | - | comma-separated bead IDs |

### Show Task

```bash
ntm --robot-bead-show=bd-abc123
```

### List Tasks

```bash
ntm --robot-beads-list --beads-status=open
ntm --robot-beads-list --beads-status=in_progress
ntm --robot-beads-list --beads-type=bug --beads-priority=1
```

### Claim Task

```bash
ntm --robot-bead-claim=bd-abc123 --bead-assignee=claude-agent
```

**Response:**
```json
{
  "success": true,
  "bead_id": "bd-abc123",
  "prev_status": "open",
  "new_status": "in_progress",
  "claimed": true
}
```

### Close Task

```bash
ntm --robot-bead-close=bd-abc123 --bead-close-reason="Fixed in commit abc123"
```

---

## Agent Control

### Send Message

```bash
# Send to all Claude agents
ntm --robot-send=myproject --msg="Fix the auth bug" --type=claude

# Send to specific panes
ntm --robot-send=myproject --msg="Run tests" --panes=1,2

# Send to all panes including user
ntm --robot-send=myproject --msg="echo hello" --all

# Dry run (preview)
ntm --robot-send=myproject --msg="test" --dry-run

# Send and wait for response
ntm --robot-send=myproject --msg="What changed?" --track --ack-timeout=60s

# Read message from file
ntm --robot-send=myproject --msg-file=/path/to/prompt.md
```

**Response:**
```json
{
  "success": true,
  "session": "myproject",
  "targets": ["1", "2"],
  "successful": ["1", "2"],
  "failed": [],
  "message_preview": "Fix the auth bug"
}
```

**Parameters:**
| Param | Required | Description |
|-------|----------|-------------|
| `--msg` | Yes* | Message content |
| `--msg-file` | Yes* | Read from file (*one required) |
| `--type` | No | Filter: claude, codex, gemini |
| `--panes` | No | Specific pane indices |
| `--all` | No | Include user pane |
| `--enter` | No | Send Enter after (default: true) |
| `--dry-run` | No | Preview only |
| `--track` | No | Wait for response |
| `--delay-ms` | No | Delay between sends |

### Spawn Agents

```bash
ntm --robot-spawn=myproject --spawn-cc=2 --spawn-cod=1 --spawn-wait
```

### Interrupt Agents

```bash
ntm --robot-interrupt=myproject --type=claude
```

---

## State Inspection

### Session Activity

```bash
# Get agent states (idle/busy/error)
ntm --robot-activity=myproject --activity-type=claude
```

**Response:**
```json
{
  "success": true,
  "session": "myproject",
  "agents": [
    {
      "pane_id": "1",
      "type": "claude",
      "state": "idle",
      "last_activity": "2026-01-23T07:00:00Z"
    }
  ]
}
```

### Context Usage

```bash
ntm --robot-context=myproject
```

### Recent Output

```bash
ntm --robot-tail=myproject --tail-lines=50
```

### Health Check

```bash
ntm --robot-diagnose=myproject
ntm --robot-diagnose=myproject --diagnose-fix  # Auto-fix issues
```

### Dashboard Summary

```bash
ntm --robot-dashboard
ntm --robot-dashboard --json
```

---

## CASS Integration (History Search)

```bash
# Search past conversations
ntm --robot-cass-search="authentication error" --cass-since=7d

# Get relevant context for a task
ntm --robot-cass-context="how to implement OAuth"

# Check CASS health
ntm --robot-cass-status
```

---

## Work Assignment

### Get Recommendations

```bash
ntm --robot-assign=myproject --strategy=balanced
ntm --robot-assign=myproject --strategy=speed --beads=bd-1,bd-2
```

**Strategies:**
- `balanced` - Even distribution
- `speed` - Prioritize fast completion
- `quality` - Prioritize thoroughness
- `dependency` - Respect task dependencies

### Bulk Assign

```bash
ntm --robot-bulk-assign=myproject --from-bv
```

---

## Automation Patterns

### Pattern 1: Task Creation and Assignment

```bash
#!/bin/bash
# Create task and assign to idle agent

# 1. Create bead
RESULT=$(ntm --robot-bead-create --bead-title="$1" --bead-type=task)
BEAD_ID=$(echo "$RESULT" | jq -r '.bead_id')

if [ "$(echo "$RESULT" | jq -r '.success')" != "true" ]; then
  echo "Failed to create bead: $(echo "$RESULT" | jq -r '.error')"
  exit 1
fi

# 2. Find idle agent
IDLE=$(ntm --robot-activity=myproject | jq -r '
  .agents[] | select(.state=="idle") | .pane_id
' | head -1)

# 3. Assign and send
if [ -n "$IDLE" ]; then
  ntm --robot-bead-claim=$BEAD_ID --bead-assignee=pane-$IDLE
  ntm --robot-send=myproject --panes=$IDLE --msg="Work on $BEAD_ID: $1"
  echo "Assigned $BEAD_ID to pane $IDLE"
else
  echo "No idle agents, bead $BEAD_ID queued"
fi
```

### Pattern 2: Context Monitor

```bash
#!/bin/bash
# Monitor context usage and warn when high

while true; do
  CONTEXT=$(ntm --robot-context=myproject)

  echo "$CONTEXT" | jq -r '.agents[] |
    select(.context_percent > 80) |
    "WARNING: \(.pane_id) at \(.context_percent)%"'

  sleep 60
done
```

### Pattern 3: Orchestrator Loop

```bash
#!/bin/bash
# Simple orchestrator: assign work to idle agents

SESSION="myproject"

while true; do
  # Get open beads
  BEADS=$(ntm --robot-beads-list --beads-status=open | jq -r '.beads[0].id')

  if [ -z "$BEADS" ] || [ "$BEADS" = "null" ]; then
    echo "No open beads"
    sleep 30
    continue
  fi

  # Get idle agents
  IDLE=$(ntm --robot-activity=$SESSION | jq -r '
    .agents[] | select(.state=="idle") | .pane_id
  ' | head -1)

  if [ -n "$IDLE" ]; then
    # Claim and assign
    ntm --robot-bead-claim=$BEADS --bead-assignee=pane-$IDLE
    TITLE=$(ntm --robot-bead-show=$BEADS | jq -r '.title')
    ntm --robot-send=$SESSION --panes=$IDLE --msg="Please work on $BEADS: $TITLE"
    echo "Assigned $BEADS to pane $IDLE"
  fi

  sleep 10
done
```

---

## Error Handling

### Error Response Format

```json
{
  "success": false,
  "timestamp": "2026-01-23T07:00:00Z",
  "error": "failed to create bead: title required",
  "error_code": "INVALID_FLAG",
  "hint": "Specify --bead-title='Your title'"
}
```

### Common Error Codes

| Code | Meaning |
|------|---------|
| `INVALID_FLAG` | Missing or invalid parameter |
| `SESSION_NOT_FOUND` | Session doesn't exist |
| `PANE_NOT_FOUND` | Pane doesn't exist |
| `INTERNAL_ERROR` | Unexpected error |
| `TIMEOUT` | Operation timed out |

### Checking Success

```bash
RESULT=$(ntm --robot-bead-create --bead-title="Test")
if [ "$(echo "$RESULT" | jq -r '.success')" = "true" ]; then
  BEAD_ID=$(echo "$RESULT" | jq -r '.bead_id')
  echo "Created: $BEAD_ID"
else
  echo "Error: $(echo "$RESULT" | jq -r '.error')"
  echo "Hint: $(echo "$RESULT" | jq -r '.hint')"
fi
```

---

## Prerequisites

### Required Binaries

```bash
# ntm - Named Tmux Manager
which ntm  # ~/.local/bin/ntm

# bd/br - Beads task tracker (required for bead commands)
which bd   # ~/.local/bin/bd -> br
which br   # ~/.local/bin/br

# tmux
which tmux # /usr/bin/tmux
```

### Project Setup

```bash
# Initialize beads in project directory
cd /path/to/project
bd init
```

---

## Quick Reference

| Task | Command |
|------|---------|
| Create task | `ntm --robot-bead-create --bead-title="..."` |
| List tasks | `ntm --robot-beads-list --beads-status=open` |
| Claim task | `ntm --robot-bead-claim=ID` |
| Close task | `ntm --robot-bead-close=ID --bead-close-reason="..."` |
| Send message | `ntm --robot-send=SESSION --msg="..."` |
| Check activity | `ntm --robot-activity=SESSION` |
| Check context | `ntm --robot-context=SESSION` |
| Health check | `ntm --robot-diagnose=SESSION` |
| Search history | `ntm --robot-cass-search="query"` |
| List all commands | `ntm --robot-capabilities` |
