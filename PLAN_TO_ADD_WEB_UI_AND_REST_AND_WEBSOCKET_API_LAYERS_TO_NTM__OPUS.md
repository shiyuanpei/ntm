# NTM Web Platform: REST API, WebSocket, and World-Class UI

> **A Comprehensive Plan for Transforming NTM into a Full-Stack Multi-Agent Orchestration Platform**
>
> *Integrating the Complete Agent Flywheel Ecosystem: Agent Mail, BV, UBS, CASS, CM, CAAM, SLB*

---

## Executive Summary

This document outlines a comprehensive plan to extend NTM (Named Tmux Manager) from a terminal/TUI application into a full-featured web platform. NTM is the **orchestration backbone** of the Agent Flywheel—a self-improving development cycle where AI coding agents work in parallel, coordinate via messaging, and compound their learnings over time.

The architecture introduces three new layers:

1. **REST API Layer** — A performant, well-documented HTTP API replicating 100% of CLI/TUI functionality across **all 8 flywheel tools**
2. **WebSocket Layer** — Real-time bidirectional streaming for logs, events, agent interactions, file reservations, scanner results, and memory updates
3. **Web UI Layer** — A world-class Next.js 16 / React 19 interface with Stripe-level polish, providing unified access to the entire flywheel ecosystem

The design prioritizes:
- **Flywheel acceleration** — Every feature designed to make the virtuous cycle spin faster
- **Full ecosystem integration** — Agent Mail, BV, UBS, CASS, CM, CAAM, SLB unified under one UI
- **Feature parity** — Every CLI command accessible via API (50+ robot mode commands)
- **Real-time capability** — Sub-100ms event propagation across all tools
- **Developer experience** — OpenAPI 3.1 spec with rich examples for AI agent consumption
- **Visual excellence** — Desktop and mobile-optimized UX with separate interaction paradigms

---

## Table of Contents

1. [The Agent Flywheel Philosophy](#1-the-agent-flywheel-philosophy)
2. [Research Findings](#2-research-findings)
3. [Architecture Overview](#3-architecture-overview)
4. [REST API Layer](#4-rest-api-layer)
5. [WebSocket Layer](#5-websocket-layer)
6. [Agent Mail Deep Integration](#6-agent-mail-deep-integration)
7. [Beads & BV Integration](#7-beads--bv-integration)
8. [CASS & Memory System Integration](#8-cass--memory-system-integration)
9. [UBS Scanner Integration](#9-ubs-scanner-integration)
10. [CAAM Account Management](#10-caam-account-management)
11. [SLB Safety Guardrails](#11-slb-safety-guardrails)
12. [Pipeline & Workflow Engine](#12-pipeline--workflow-engine)
13. [Web UI Layer](#13-web-ui-layer)
14. [Desktop vs Mobile UX Strategy](#14-desktop-vs-mobile-ux-strategy)
15. [Agent SDK Integration Strategy](#15-agent-sdk-integration-strategy)
16. [Implementation Phases](#16-implementation-phases)
17. [File Structure](#17-file-structure)
18. [Technical Specifications](#18-technical-specifications)

---

## 1. The Agent Flywheel Philosophy

### 1.1 What Is The Agent Flywheel?

The Agent Flywheel is a **self-improving development cycle** where:

```
┌─────────────────────────────────────────────────────────────────┐
│                    THE AGENT FLYWHEEL                           │
│                                                                 │
│         ┌─────────┐                                            │
│         │  PLAN   │◄────────────────────────────────┐          │
│         │  (BV)   │                                 │          │
│         └────┬────┘                                 │          │
│              │                                      │          │
│              ▼                                      │          │
│         ┌─────────┐                                 │          │
│         │COORDINATE                                 │          │
│         │(Agent   │                                 │          │
│         │ Mail)   │                                 │          │
│         └────┬────┘                                 │          │
│              │                                      │          │
│              ▼                                      │          │
│         ┌─────────┐         ┌─────────┐            │          │
│         │ EXECUTE │────────▶│  SCAN   │            │          │
│         │ (NTM +  │         │  (UBS)  │            │          │
│         │ Agents) │         └────┬────┘            │          │
│         └─────────┘              │                 │          │
│                                  ▼                 │          │
│                             ┌─────────┐            │          │
│                             │REMEMBER │────────────┘          │
│                             │(CASS+CM)│                       │
│                             └─────────┘                       │
│                                                                 │
│  Each cycle is better than the last because:                   │
│  • Memory improves (CM gets smarter)                           │
│  • Sessions are searchable (find past solutions)               │
│  • Agents coordinate (no duplicated work)                      │
│  • Quality gates enforce standards (UBS)                       │
│  • Context is preserved (Agent Mail + CM)                      │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 The Eight Tools of the Flywheel

| # | Tool | Purpose | Integration Priority |
|---|------|---------|---------------------|
| 1 | **NTM** | Session orchestration & agent spawning | Core (existing) |
| 2 | **Agent Mail** | Agent messaging & file coordination | Critical |
| 3 | **BV** | Task management & graph analysis | Critical |
| 4 | **UBS** | Code quality scanning | High |
| 5 | **CASS** | Session history search & indexing | High |
| 6 | **CM** | Procedural memory for agents | High |
| 7 | **CAAM** | Authentication credential rotation | Medium |
| 8 | **SLB** | Safety guardrails (two-person rule) | Medium |

### 1.3 How The Web UI Accelerates The Flywheel

The web UI transforms each phase:

| Phase | CLI Experience | Web UI Experience |
|-------|----------------|-------------------|
| **PLAN** | `bv` TUI, `bd ready` | Visual Kanban, dependency graph, drag-drop prioritization |
| **COORDINATE** | `am` commands, inbox polling | Real-time chat, file reservation map, @mentions |
| **EXECUTE** | `ntm spawn`, tmux attach | Visual agent grid, live terminals, one-click spawn |
| **SCAN** | `ubs .` output | Dashboard with severity charts, inline annotations |
| **REMEMBER** | `cm context`, `cass search` | Semantic search UI, memory timeline, rule browser |

### 1.4 Design Principle: Flywheel-First

Every feature should answer: **"Does this make the flywheel spin faster?"**

- ✅ Real-time file reservation map → Prevents conflicts, faster coordination
- ✅ Visual dependency graph → Better prioritization, faster planning
- ✅ Inline UBS annotations → Faster bug fixing, better quality
- ✅ Memory search UI → Faster context retrieval, better first attempts
- ❌ Pretty animations with no function → Slower page loads, distraction

---

## 2. Research Findings

### 2.1 Agent Client Protocol (ACP)

The [Agent Client Protocol](https://agentclientprotocol.com/) is an emerging open standard (Apache 2.0) created by Zed for connecting AI coding agents to editors/IDEs. Key findings:

- **JSON-RPC 2.0 based** — Bidirectional communication over stdio
- **Industry adoption** — JetBrains, Neovim, and Google (Gemini CLI reference implementation)
- **Complements MCP** — MCP handles data/tools; ACP handles agent-editor integration
- **Adapters available** — Open-source adapters for Claude Code, Codex, Gemini CLI

**Implication for NTM:** ACP provides a standardized way to communicate with agents that could eventually supplement or replace tmux-based text streaming. We should design the API to support both paradigms.

### 2.2 Official Agent SDKs

| SDK | Package | Version | Key Features |
|-----|---------|---------|--------------|
| **Claude Agent SDK** | `@anthropic-ai/claude-agent-sdk` | 0.1.76 | Async streaming, tool use, file ops |
| **OpenAI Codex SDK** | `@openai/codex-sdk` | Latest | JSONL events over stdin/stdout, thread persistence |
| **Google GenAI SDK** | `@google/genai` | 1.34.0 (GA) | 1M token context, MCP support |

**Implication:** We can offer a "direct SDK mode" as an alternative to tmux spawning, giving users choice between:
- **Tmux mode** — Current approach, battle-tested, visual terminal access
- **SDK mode** — Lower overhead, programmatic control, no tmux dependency

### 2.3 Next.js 16 / React 19.2

Released October 2025, Next.js 16 brings:

- **Turbopack stable** — 10× faster Fast Refresh (default in dev and build)
- **React Compiler 1.0** — Automatic memoization, zero manual optimization
- **React 19.2 features**:
  - `View Transitions` — Native animation between route changes
  - `Activity` — Background rendering with state preservation
  - `useEffectEvent` — Non-reactive Effect logic extraction
- **Cache Components** — Explicit `"use cache"` directive (opt-in caching)
- **Enhanced routing** — Layout deduplication, incremental prefetching

### 2.4 MCP Agent Mail Protocol

NTM's existing Agent Mail integration uses HTTP JSON-RPC to `localhost:8765`. Key capabilities:

- **Project & Agent Management**: `EnsureProject`, `RegisterAgent`, `CreateAgentIdentity`
- **Messaging**: `SendMessage`, `ReplyMessage`, `FetchInbox`, `SearchMessages`, `SummarizeThread`
- **File Reservations**: `ReservePaths`, `ReleaseReservations`, `RenewReservations`, `ForceReleaseReservation`
- **Contact Management**: `RequestContact`, `RespondContact`, `ListContacts`
- **Macros**: `StartSession`, `PrepareThread`, `ContactHandshake`
- **Overseer Mode**: `SendOverseerMessage` (bypass contact policies)
- **Pre-commit Guards**: `InstallPrecommitGuard`, `UninstallPrecommitGuard`

### 2.5 BV Robot Mode Commands

NTM integrates with BV for task management:

- `GetTriage()` — Comprehensive triage with scoring, recommendations, quick wins (30s cache)
- `GetInsights()` — Graph analysis: bottlenecks, keystones, hubs, authorities, cycles
- `GetPriority()` — Priority recommendations with impact scoring
- `GetPlan()` — Parallel execution plan for work distribution

### 2.6 CASS Search Capabilities

- `Search(query, options)` — Full-text search with filters (agent, workspace, time range)
- `SearchForContext(query, workspace)` — Context retrieval for task planning
- Aggregations by agent, workspace, and tags
- Pagination with cursor support

### 2.7 Real-Time Streaming Best Practices

From 2025 WebSocket research:

- **Bidirectional necessity** — Terminal interaction requires full-duplex
- **Reconnection handling** — Must be application-specific with state recovery
- **Horizontal scaling** — Redis adapter pattern for multi-server broadcast
- **Message ordering** — Critical for terminal output coherence
- **Edge deployment** — Reduce latency via geo-distributed WebSocket servers

### 2.8 TanStack Query + WebSocket Pattern

TanStack Query v5 doesn't have first-class WebSocket support, but the recommended pattern:

```typescript
// Initial fetch with useQuery
const { data } = useQuery({ queryKey: ['session', id], queryFn: fetchSession });

// WebSocket updates via queryClient.setQueryData
ws.onmessage = (event) => {
  queryClient.setQueryData(['session', id], (old) => merge(old, event.data));
};
```

The new `streamedQuery` API in v5 provides 3× faster perceived performance for streaming data.

---

## 3. Architecture Overview

### 3.1 High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                           NTM WEB PLATFORM                                        │
│                    (Agent Flywheel Command Center)                                │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │                         WEB UI (Next.js 16)                                 │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐         │ │
│  │  │Dashboard │ │ Sessions │ │  Beads   │ │  Memory  │ │ Scanner  │         │ │
│  │  │  View    │ │  View    │ │  Board   │ │  Search  │ │Dashboard │         │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘         │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐         │ │
│  │  │ Palette  │ │  Mail    │ │ Accounts │ │ Pipeline │ │  Mobile  │         │ │
│  │  │  View    │ │  Inbox   │ │  Manager │ │  Builder │ │  Views   │         │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘         │ │
│  │                                                                            │ │
│  │  ┌──────────────────────────────────────────────────────────────────────┐ │ │
│  │  │           TanStack Query + WebSocket Provider + Zustand              │ │ │
│  │  └──────────────────────────────────────────────────────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                              │                    │                              │
│                         HTTP/REST            WebSocket                           │
│                              │                    │                              │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │                      GO HTTP SERVER (net/http + chi)                        │ │
│  │  ┌────────────────────────────┐  ┌────────────────────────────────────┐    │ │
│  │  │       REST ROUTER          │  │        WEBSOCKET HUB               │    │ │
│  │  │                            │  │                                    │    │ │
│  │  │  /api/v1/sessions          │  │  Channels:                         │    │ │
│  │  │  /api/v1/agents            │  │  • output    (pane streams)        │    │ │
│  │  │  /api/v1/beads             │  │  • status    (agent states)        │    │ │
│  │  │  /api/v1/mail              │  │  • health    (health events)       │    │ │
│  │  │  /api/v1/reservations      │  │  • alerts    (notifications)       │    │ │
│  │  │  /api/v1/cass              │  │  • beads     (task updates)        │    │ │
│  │  │  /api/v1/memory            │  │  • mail      (messages)            │    │ │
│  │  │  /api/v1/scanner           │  │  • files     (changes)             │    │ │
│  │  │  /api/v1/accounts          │  │  • scanner   (scan results)        │    │ │
│  │  │  /api/v1/pipelines         │  │  • pipeline  (workflow events)     │    │ │
│  │  │  /api/v1/safety            │  │  • memory    (CM updates)          │    │ │
│  │  └────────────────────────────┘  └────────────────────────────────────┘    │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                       │                                          │
│  ┌────────────────────────────────────┴───────────────────────────────────────┐ │
│  │                        NTM CORE (Existing Go Packages)                      │ │
│  │                                                                             │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │ │
│  │  │  tmux/  │ │ robot/  │ │ config/ │ │ agents/ │ │   bv/   │ │  cass/  │  │ │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │ │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │ │
│  │  │agentmail│ │ scanner │ │checkpoint│ │palette/ │ │pipeline │ │resilience│  │ │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                       │                                          │
│                    ┌──────────────────┼──────────────────┐                       │
│                    │                  │                  │                       │
│  ┌─────────────────▼───┐  ┌──────────▼──────────┐  ┌────▼────────────────────┐  │
│  │    TMUX SERVER      │  │   AGENT MAIL MCP    │  │   EXTERNAL TOOLS        │  │
│  │  (Sessions/Panes)   │  │   (localhost:8765)  │  │  • UBS (scanner)        │  │
│  └─────────────────────┘  └─────────────────────┘  │  • CASS (search)        │  │
│                                                     │  • CM (memory)          │  │
│  ┌─────────────────────────────────────────────┐   │  • CAAM (accounts)      │  │
│  │           AI CODING AGENTS                   │   │  • SLB (safety)         │  │
│  │  Claude Code │ Codex CLI │ Gemini CLI        │   └─────────────────────────┘  │
│  └─────────────────────────────────────────────┘                                 │
│                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Design Principles

1. **Zero functionality loss** — Every CLI command has an API equivalent
2. **Robot mode as foundation** — REST responses mirror existing `--robot-*` JSON structures
3. **Flywheel-first** — Every feature accelerates the virtuous cycle
4. **Streaming-first** — WebSocket for all real-time data; REST for commands/queries
5. **Unified ecosystem** — All 8 tools accessible from single UI
6. **Layered abstraction** — API layer is thin; business logic stays in existing packages
7. **Backward compatible** — CLI continues to work unchanged
8. **Progressive enhancement** — Web UI enhances but doesn't replace terminal workflow

### 3.3 Technology Decisions

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **REST Server** | Go `net/http` + `chi` router | Native Go, performant, middleware ecosystem |
| **WebSocket** | `gorilla/websocket` | Battle-tested, concurrent-safe, ping/pong |
| **Event Bus** | Internal Go pub/sub | Already exists in NTM, 100-event ring buffer |
| **API Docs** | OpenAPI 3.1 + Swagger UI | Industry standard, code generation |
| **Frontend** | Next.js 16 + React 19 | Latest features, Turbopack, React Compiler |
| **State** | TanStack Query v5 + Zustand | Server state + client state separation |
| **Styling** | Tailwind CSS 4 + Framer Motion | Utility-first, animation primitives |
| **Terminal** | xterm.js | Full terminal emulation in browser |
| **Icons** | Lucide React | Consistent, tree-shakeable |
| **Charts** | Recharts + Tremor | Dashboard visualizations |
| **Graphs** | React Flow | Dependency graph visualization |

---

## 4. REST API Layer

### 4.1 API Design Philosophy

The REST API follows these principles:

1. **Resource-oriented** — Sessions, agents, panes, beads, reservations as resources
2. **Consistent responses** — All responses follow the robot mode structure
3. **Idempotent where possible** — PUT/DELETE operations are idempotent
4. **Rich error responses** — Error codes, messages, and actionable hints
5. **AI-agent friendly** — Comprehensive examples for LLM consumption

### 4.2 Base URL Structure

```
Production:  https://api.ntm.local/v1
Development: http://localhost:8080/api/v1
```

### 4.3 Authentication

```yaml
# API Key authentication (header)
Authorization: Bearer ntm_sk_live_xxxxxxxxxxxx

# Or query parameter for WebSocket connections
?api_key=ntm_sk_live_xxxxxxxxxxxx
```

### 4.4 Standard Response Envelope

All responses follow this structure (matching robot mode):

```typescript
interface ApiResponse<T> {
  success: boolean;
  timestamp: string;        // RFC3339 UTC
  data?: T;                 // On success
  error?: string;           // Human-readable error
  error_code?: string;      // Programmatic error code
  hint?: string;            // Actionable guidance
  _agent_hints?: {          // For AI agent consumers
    summary: string;
    suggested_actions: Action[];
    warnings: string[];
  };
}
```

### 4.5 Complete Endpoint Catalog

#### 4.5.1 Sessions (`/api/v1/sessions`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/sessions` | `--robot-status` | List all sessions |
| `POST` | `/sessions` | `ntm create` | Create empty session |
| `GET` | `/sessions/{name}` | `--robot-status` | Get session details |
| `DELETE` | `/sessions/{name}` | `ntm kill` | Kill session |
| `POST` | `/sessions/{name}/spawn` | `--robot-spawn` | Create with agents |
| `POST` | `/sessions/{name}/attach` | `ntm attach` | Mark attached |
| `GET` | `/sessions/{name}/snapshot` | `--robot-snapshot` | Full state capture |

#### 4.5.2 Agents & Panes (`/api/v1/sessions/{name}/...`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/sessions/{name}/panes` | `--robot-status` | List panes |
| `GET` | `/sessions/{name}/panes/{idx}` | `--robot-inspect-pane` | Pane details |
| `POST` | `/sessions/{name}/agents` | `ntm add` | Add agents |
| `GET` | `/sessions/{name}/context` | `--robot-context` | Context usage |
| `GET` | `/sessions/{name}/health` | `--robot-health` | Health status |
| `GET` | `/sessions/{name}/activity` | `--robot-activity` | Activity states |
| `GET` | `/sessions/{name}/metrics` | `--robot-metrics` | Session metrics |

#### 4.5.3 Messaging & Control (`/api/v1/sessions/{name}/...`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `POST` | `/sessions/{name}/send` | `--robot-send` | Send prompt |
| `POST` | `/sessions/{name}/interrupt` | `--robot-interrupt` | Send Ctrl+C |
| `POST` | `/sessions/{name}/wait` | `--robot-wait` | Wait for condition |
| `GET` | `/sessions/{name}/route` | `--robot-route` | Routing recommendation |
| `GET` | `/sessions/{name}/output/tail` | `--robot-tail` | Recent output |
| `GET` | `/sessions/{name}/history` | `--robot-history` | Command history |
| `POST` | `/sessions/{name}/replay` | `--robot-replay` | Replay command |

#### 4.5.4 Checkpoints (`/api/v1/sessions/{name}/checkpoints`)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/sessions/{name}/checkpoints` | List checkpoints |
| `POST` | `/sessions/{name}/checkpoints` | Create checkpoint |
| `GET` | `/sessions/{name}/checkpoints/{id}` | Get checkpoint |
| `DELETE` | `/sessions/{name}/checkpoints/{id}` | Delete checkpoint |
| `POST` | `/sessions/{name}/checkpoints/{id}/restore` | Restore |

#### 4.5.5 Files & Conflicts (`/api/v1/sessions/{name}/files`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/sessions/{name}/files` | `--robot-files` | File changes |
| `GET` | `/sessions/{name}/files/conflicts` | `--robot-diff` | Conflict detection |

#### 4.5.6 Beads (`/api/v1/beads`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/beads` | `--robot-beads-list` | List beads |
| `POST` | `/beads` | `--robot-bead-create` | Create bead |
| `GET` | `/beads/{id}` | `--robot-bead-show` | Get bead |
| `PATCH` | `/beads/{id}` | `bd update` | Update bead |
| `POST` | `/beads/{id}/claim` | `--robot-bead-claim` | Claim bead |
| `POST` | `/beads/{id}/close` | `--robot-bead-close` | Close bead |
| `GET` | `/beads/ready` | `bd ready` | Ready beads |
| `GET` | `/beads/triage` | `bv --robot-triage` | Triage analysis |
| `GET` | `/beads/plan` | `--robot-plan` | Execution plan |
| `GET` | `/beads/insights` | `bv --robot-insights` | Graph insights |

#### 4.5.7 Agent Mail (`/api/v1/mail`)

| Method | Endpoint | MCP Equivalent | Description |
|--------|----------|----------------|-------------|
| `GET` | `/mail/inbox/{agent}` | `FetchInbox` | Agent inbox |
| `POST` | `/mail/send` | `SendMessage` | Send message |
| `POST` | `/mail/reply` | `ReplyMessage` | Reply to message |
| `POST` | `/mail/{id}/read` | `MarkMessageRead` | Mark read |
| `POST` | `/mail/{id}/ack` | `AcknowledgeMessage` | Acknowledge |
| `GET` | `/mail/threads/{id}` | `SummarizeThread` | Thread summary |
| `GET` | `/mail/search` | `SearchMessages` | Search messages |
| `POST` | `/mail/overseer` | `SendOverseerMessage` | Overseer message |

#### 4.5.8 File Reservations (`/api/v1/reservations`)

| Method | Endpoint | MCP Equivalent | Description |
|--------|----------|----------------|-------------|
| `GET` | `/reservations` | `ListReservations` | All reservations |
| `GET` | `/reservations/{agent}` | `ListReservations` | Agent reservations |
| `POST` | `/reservations` | `ReservePaths` | Create reservation |
| `DELETE` | `/reservations/{id}` | `ReleaseReservations` | Release |
| `POST` | `/reservations/{id}/renew` | `RenewReservations` | Extend TTL |
| `POST` | `/reservations/{id}/force-release` | `ForceReleaseReservation` | Force release |

#### 4.5.9 CASS & Memory (`/api/v1/cass`, `/api/v1/memory`)

| Method | Endpoint | CLI Equivalent | Description |
|--------|----------|----------------|-------------|
| `GET` | `/cass/status` | `--robot-cass-status` | CASS health |
| `GET` | `/cass/search` | `--robot-cass-search` | Search sessions |
| `GET` | `/cass/context` | `--robot-cass-context` | Get context |
| `GET` | `/cass/insights` | `--robot-cass-insights` | Aggregated insights |
| `GET` | `/memory/context` | `cm context` | Memory context |
| `POST` | `/memory/reflect` | `cm reflect` | Update memory |
| `GET` | `/memory/playbook` | `cm playbook list` | List rules |

#### 4.5.10 Scanner (`/api/v1/scanner`)

| Method | Endpoint | CLI Equivalent | Description |
|--------|----------|----------------|-------------|
| `POST` | `/scanner/scan` | `ubs .` | Run scan |
| `GET` | `/scanner/results/{id}` | — | Get scan results |
| `GET` | `/scanner/latest` | — | Latest scan |
| `GET` | `/scanner/status` | `ubs doctor` | Scanner status |
| `POST` | `/scanner/watch/start` | Auto-scanner | Start watching |
| `POST` | `/scanner/watch/stop` | — | Stop watching |

#### 4.5.11 Accounts (`/api/v1/accounts`)

| Method | Endpoint | CLI Equivalent | Description |
|--------|----------|----------------|-------------|
| `GET` | `/accounts` | `caam status` | List accounts |
| `GET` | `/accounts/{type}` | — | Accounts by type |
| `POST` | `/accounts/{type}/activate` | `caam activate` | Activate account |
| `GET` | `/accounts/active` | — | Active accounts |

#### 4.5.12 Pipelines (`/api/v1/pipelines`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/pipelines` | `--robot-pipeline-list` | List pipelines |
| `POST` | `/pipelines/run` | `--robot-pipeline-run` | Run pipeline |
| `GET` | `/pipelines/{id}` | `--robot-pipeline` | Get status |
| `POST` | `/pipelines/{id}/cancel` | `--robot-pipeline-cancel` | Cancel |

#### 4.5.13 Alerts & System (`/api/v1/...`)

| Method | Endpoint | Robot Equivalent | Description |
|--------|----------|------------------|-------------|
| `GET` | `/alerts` | `--robot-alerts` | Active alerts |
| `POST` | `/alerts/{id}/dismiss` | `--robot-dismiss-alert` | Dismiss |
| `GET` | `/config` | `ntm config show` | Configuration |
| `GET` | `/config/palette` | `--robot-palette` | Palette commands |
| `GET` | `/config/recipes` | `--robot-recipes` | Spawn recipes |
| `GET` | `/health` | — | System health |
| `GET` | `/version` | `--robot-version` | Version info |

### 4.6 Example: Create Session with Agents

```http
POST /api/v1/sessions/myproject/spawn
Content-Type: application/json

{
  "agents": {
    "claude": 3,
    "codex": 2,
    "gemini": 1
  },
  "options": {
    "auto_restart": true,
    "cass_context_query": "authentication patterns",
    "stagger_ms": 500,
    "recipe": "full-stack",
    "register_agent_mail": true
  }
}
```

Response:
```json
{
  "success": true,
  "timestamp": "2025-01-07T15:30:00Z",
  "data": {
    "session": "myproject",
    "created": true,
    "agents": [
      {"pane": "myproject__cc_1", "type": "claude", "index": 1, "agent_mail_id": "GreenCastle"},
      {"pane": "myproject__cc_2", "type": "claude", "index": 2, "agent_mail_id": "BlueLake"},
      {"pane": "myproject__cc_3", "type": "claude", "index": 3, "agent_mail_id": "RedStone"}
    ],
    "total_panes": 7,
    "agent_mail": {
      "project_registered": true,
      "agents_registered": 6
    }
  },
  "_agent_hints": {
    "summary": "Created session with 6 agents + 1 user pane, all registered with Agent Mail",
    "suggested_actions": [
      {"action": "reserve_files", "endpoint": "POST /reservations"},
      {"action": "send_prompt", "endpoint": "POST /sessions/myproject/send"}
    ]
  }
}
```

---

## 5. WebSocket Layer

### 5.1 Connection Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WebSocket Connection Manager                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    CHANNEL ROUTER                            │   │
│  │                                                              │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │ output   │ │ status   │ │ health   │ │ alerts   │       │   │
│  │  │(terminal)│ │(agents)  │ │(system)  │ │(notifs)  │       │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │ beads    │ │ mail     │ │ files    │ │ scanner  │       │   │
│  │  │(tasks)   │ │(messages)│ │(changes) │ │(results) │       │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │   │
│  │  │pipeline  │ │ memory   │ │reserv.   │ │ accounts │       │   │
│  │  │(workflow)│ │(CM)      │ │(locks)   │ │(rotation)│       │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │   │
│  │                                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                              │                                      │
│                              ▼                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    EVENT SOURCES                             │   │
│  │                                                              │   │
│  │  • Tmux pane output polling (configurable interval)         │   │
│  │  • Agent Mail inbox polling                                  │   │
│  │  • BV triage cache invalidation                             │   │
│  │  • UBS auto-scanner results                                  │   │
│  │  • CASS index updates                                        │   │
│  │  • CM memory changes                                         │   │
│  │  • File system watchers                                      │   │
│  │  • Health check results                                      │   │
│  │  • Pipeline state changes                                    │   │
│  │  • Account rotation events                                   │   │
│  │                                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.2 Connection Endpoint

```
WebSocket URL: wss://api.ntm.local/v1/ws
              ws://localhost:8080/api/v1/ws (development)

Query Parameters:
  - api_key: Authentication token
  - session: Optional session filter
  - channels: Comma-separated channel list
```

### 5.3 Message Protocol

```typescript
interface WebSocketMessage {
  id: string;              // Unique message ID
  type: MessageType;       // Message type enum
  channel: Channel;        // Which channel
  timestamp: string;       // RFC3339 UTC
  payload: unknown;        // Varies by type
  request_id?: string;     // For request/response correlation
}

type MessageType =
  | 'subscribe'      // Client → Server
  | 'unsubscribe'    // Client → Server
  | 'subscribed'     // Server → Client
  | 'event'          // Server → Client
  | 'output'         // Server → Client (terminal)
  | 'error'          // Server → Client
  | 'ping' | 'pong'  // Bidirectional
  | 'command'        // Client → Server
  | 'result';        // Server → Client

type Channel =
  | 'output'         // Pane output streams
  | 'status'         // Agent status changes
  | 'health'         // Health events
  | 'alerts'         // Alert notifications
  | 'beads'          // Bead/task updates
  | 'mail'           // Agent Mail messages
  | 'files'          // File change events
  | 'reservations'   // File reservation changes
  | 'scanner'        // UBS scan results
  | 'pipeline'       // Pipeline workflow events
  | 'memory'         // CM memory updates
  | 'accounts'       // CAAM account changes
  | 'system';        // System-wide events
```

### 5.4 Channel Specifications

#### 5.4.1 Output Channel (`output`)

Streams real-time pane output.

```json
{
  "type": "output",
  "channel": "output",
  "timestamp": "2025-01-07T15:30:00.123Z",
  "payload": {
    "session": "myproject",
    "pane": "myproject__cc_1",
    "pane_index": 1,
    "agent_type": "claude",
    "agent_mail_name": "GreenCastle",
    "content": "Analyzing the authentication module...\n",
    "sequence": 12345,
    "is_error": false,
    "detected_state": "working"
  }
}
```

#### 5.4.2 Beads Channel (`beads`)

Task/issue updates from BV.

```json
{
  "type": "event",
  "channel": "beads",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "bead_updated",
    "bead_id": "bd-45",
    "title": "Implement user authentication",
    "changes": {
      "status": {"from": "open", "to": "in_progress"},
      "assignee": {"from": null, "to": "GreenCastle"}
    },
    "graph_impact": {
      "unblocks": ["bd-48", "bd-49"],
      "is_bottleneck": true,
      "bottleneck_score": 0.85
    }
  }
}
```

#### 5.4.3 Mail Channel (`mail`)

Agent Mail message notifications.

```json
{
  "type": "event",
  "channel": "mail",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "message_received",
    "message_id": 1234,
    "from": "GreenCastle",
    "to": ["BlueLake", "RedStone"],
    "subject": "Auth module ready for review",
    "importance": "normal",
    "ack_required": true,
    "thread_id": "FEAT-auth-123",
    "preview": "I've completed the JWT validation..."
  }
}
```

#### 5.4.4 Reservations Channel (`reservations`)

File reservation changes.

```json
{
  "type": "event",
  "channel": "reservations",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "reservation_granted",
    "reservation_id": 42,
    "agent": "GreenCastle",
    "paths": ["src/auth/**/*.go"],
    "exclusive": true,
    "expires_at": "2025-01-07T16:30:00Z",
    "conflicts_with": []
  }
}
```

#### 5.4.5 Scanner Channel (`scanner`)

UBS scan results.

```json
{
  "type": "event",
  "channel": "scanner",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "scan_complete",
    "scan_id": "scan_abc123",
    "trigger": "auto",
    "project": "myproject",
    "duration_ms": 3200,
    "totals": {
      "critical": 0,
      "warning": 3,
      "info": 12,
      "files": 45
    },
    "findings": [
      {
        "file": "src/auth/jwt.go",
        "line": 42,
        "severity": "warning",
        "message": "Unused variable 'tokenExpiry'",
        "suggestion": "Remove or use the variable"
      }
    ],
    "gate_passed": true
  }
}
```

#### 5.4.6 Memory Channel (`memory`)

CM memory updates.

```json
{
  "type": "event",
  "channel": "memory",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "memory_updated",
    "source": "reflection",
    "new_rules": 3,
    "updated_rules": 1,
    "categories_affected": ["authentication", "error-handling"],
    "sample_rule": {
      "id": "rule_xyz",
      "content": "Always validate JWT expiry before trusting claims",
      "confidence": 0.92
    }
  }
}
```

#### 5.4.7 Pipeline Channel (`pipeline`)

Workflow execution events.

```json
{
  "type": "event",
  "channel": "pipeline",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "step_completed",
    "pipeline_id": "pipe_123",
    "pipeline_name": "auth-feature",
    "step": {
      "index": 2,
      "name": "run_tests",
      "status": "success",
      "duration_ms": 45000,
      "output_summary": "42 tests passed"
    },
    "progress": {
      "completed": 2,
      "total": 5,
      "percent": 40
    }
  }
}
```

#### 5.4.8 Accounts Channel (`accounts`)

CAAM account rotation events.

```json
{
  "type": "event",
  "channel": "accounts",
  "timestamp": "2025-01-07T15:30:00Z",
  "payload": {
    "event": "account_rotated",
    "agent_type": "claude",
    "previous_account": "primary",
    "new_account": "backup-1",
    "reason": "rate_limit_hit",
    "cooldown_until": "2025-01-07T15:45:00Z"
  }
}
```

---

## 6. Agent Mail Deep Integration

### 6.1 Agent Mail Architecture in Web UI

Agent Mail is the **coordination backbone** of the flywheel. The web UI provides:

1. **Real-time Inbox** — Live message updates via WebSocket
2. **Thread View** — Conversation threading with summaries
3. **File Reservation Map** — Visual representation of who owns what
4. **Contact Management** — Approve/deny agent contact requests
5. **Overseer Panel** — Human oversight with priority messaging

### 6.2 File Reservation Visualization

```
┌─────────────────────────────────────────────────────────────────┐
│                    FILE RESERVATION MAP                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Project: myproject                    Total Files: 127         │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  src/                                                    │   │
│  │  ├── auth/          🟣 GreenCastle (exclusive, 45m)     │   │
│  │  │   ├── jwt.go     🟣                                  │   │
│  │  │   ├── session.go 🟣                                  │   │
│  │  │   └── oauth.go   🟣                                  │   │
│  │  ├── api/           🔵 BlueLake (exclusive, 30m)        │   │
│  │  │   ├── users.go   🔵                                  │   │
│  │  │   └── routes.go  🔵                                  │   │
│  │  └── models/        ⚪ Available                         │   │
│  │      └── user.go    ⚪                                   │   │
│  │                                                          │   │
│  │  tests/             🟡 RedStone (shared, 60m)           │   │
│  │  ├── auth_test.go   🟡                                  │   │
│  │  └── api_test.go    🟡                                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Legend: 🟣 Claude  🔵 Codex  🟡 Gemini  ⚪ Available          │
│                                                                 │
│  [Request Reservation]  [View Conflicts]  [Force Release]       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.3 Agent Mail REST Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/mail/projects` | Ensure project exists |
| `POST` | `/mail/agents` | Register agent identity |
| `GET` | `/mail/agents/{project}` | List project agents |
| `GET` | `/mail/inbox/{agent}` | Fetch inbox |
| `POST` | `/mail/send` | Send message |
| `POST` | `/mail/reply/{id}` | Reply to message |
| `POST` | `/mail/{id}/read` | Mark as read |
| `POST` | `/mail/{id}/ack` | Acknowledge |
| `GET` | `/mail/threads/{id}` | Get thread |
| `GET` | `/mail/threads/{id}/summary` | Summarize thread |
| `GET` | `/mail/search` | Search messages |
| `POST` | `/mail/overseer` | Send overseer message |
| `GET` | `/reservations` | List all reservations |
| `POST` | `/reservations` | Create reservation |
| `DELETE` | `/reservations/{id}` | Release |
| `POST` | `/reservations/{id}/renew` | Extend TTL |
| `POST` | `/reservations/{id}/force` | Force release |
| `GET` | `/contacts/{agent}` | List contacts |
| `POST` | `/contacts/request` | Request contact |
| `POST` | `/contacts/respond` | Accept/deny |
| `POST` | `/guards/install` | Install pre-commit |
| `POST` | `/guards/uninstall` | Uninstall |

### 6.4 Agent Mail UI Components

```tsx
// components/mail/inbox.tsx
interface InboxProps {
  agent: string;
  projectKey: string;
}

export function AgentInbox({ agent, projectKey }: InboxProps) {
  const { data: messages } = useInbox(agent, projectKey);
  const { data: unreadCount } = useUnreadCount(agent);

  return (
    <div className="flex flex-col h-full">
      {/* Header with unread badge */}
      <InboxHeader agent={agent} unreadCount={unreadCount} />

      {/* Message list */}
      <ScrollArea className="flex-1">
        <AnimatePresence>
          {messages?.map((msg) => (
            <MessageRow
              key={msg.id}
              message={msg}
              onRead={() => markRead(msg.id)}
              onAck={() => acknowledge(msg.id)}
            />
          ))}
        </AnimatePresence>
      </ScrollArea>

      {/* Compose button */}
      <ComposeButton onClick={() => openCompose(agent)} />
    </div>
  );
}

// components/mail/reservation-map.tsx
export function ReservationMap({ projectKey }: { projectKey: string }) {
  const { data: reservations } = useReservations(projectKey);
  const { data: fileTree } = useFileTree(projectKey);

  return (
    <div className="h-full">
      <FileTreeView
        tree={fileTree}
        reservations={reservations}
        onRequestReservation={handleRequest}
        onReleaseReservation={handleRelease}
      />
    </div>
  );
}
```

---

## 7. Beads & BV Integration

### 7.1 Beads Board (Kanban View)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           BEADS BOARD                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  [+ New Bead]  [Triage View]  [Graph View]  [Filter: All ▾]  [Sort: Score ▾]│
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐                │
│  │    READY (9)   │  │ IN PROGRESS(3) │  │   BLOCKED (2)  │                │
│  ├────────────────┤  ├────────────────┤  ├────────────────┤                │
│  │                │  │                │  │                │                │
│  │ ┌────────────┐ │  │ ┌────────────┐ │  │ ┌────────────┐ │                │
│  │ │ bd-45      │ │  │ │ bd-43      │ │  │ │ bd-41      │ │                │
│  │ │ Auth API   │ │  │ │ DB Layer   │ │  │ │ CI Setup   │ │                │
│  │ │ ★★★★☆ 0.92│ │  │ │ 🟣 Green   │ │  │ │ Blocked by │ │                │
│  │ │ 🔓 Unblocks│ │  │ │ Castle     │ │  │ │ bd-40      │ │                │
│  │ │   3 tasks  │ │  │ └────────────┘ │  │ └────────────┘ │                │
│  │ └────────────┘ │  │                │  │                │                │
│  │                │  │ ┌────────────┐ │  │ ┌────────────┐ │                │
│  │ ┌────────────┐ │  │ │ bd-47      │ │  │ │ bd-42      │ │                │
│  │ │ bd-51      │ │  │ │ API Tests  │ │  │ │ Deploy     │ │                │
│  │ │ Fix README │ │  │ │ 🔵 Blue    │ │  │ │ Blocked by │ │                │
│  │ │ Quick Win  │ │  │ │ Lake       │ │  │ │ bd-41      │ │                │
│  │ │ ⚡ Low     │ │  │ └────────────┘ │  │ └────────────┘ │                │
│  │ │   effort   │ │  │                │  │                │                │
│  │ └────────────┘ │  │                │  │                │                │
│  │                │  │                │  │                │                │
│  │ [+2 more...]  │  │                │  │                │                │
│  │                │  │                │  │                │                │
│  └────────────────┘  └────────────────┘  └────────────────┘                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Dependency Graph View

Using React Flow for interactive graph visualization:

```tsx
// components/beads/dependency-graph.tsx
export function DependencyGraph({ beads }: { beads: Bead[] }) {
  const nodes = beads.map(bead => ({
    id: bead.id,
    type: 'beadNode',
    data: {
      bead,
      isBottleneck: bead.graphPosition?.isBottleneck,
      isKeystone: bead.graphPosition?.isKeystone,
    },
    position: calculatePosition(bead),
  }));

  const edges = beads.flatMap(bead =>
    bead.blockedBy?.map(dep => ({
      id: `${dep}-${bead.id}`,
      source: dep,
      target: bead.id,
      animated: bead.status === 'blocked',
    })) ?? []
  );

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={{ beadNode: BeadNode }}
      fitView
    >
      <Background />
      <Controls />
      <MiniMap />
    </ReactFlow>
  );
}

// Custom node component
function BeadNode({ data }: NodeProps) {
  const { bead, isBottleneck, isKeystone } = data;

  return (
    <div className={cn(
      "p-3 rounded-lg border-2",
      isBottleneck && "border-red-500 bg-red-500/10",
      isKeystone && "border-yellow-500 bg-yellow-500/10",
      !isBottleneck && !isKeystone && "border-surface1 bg-surface0"
    )}>
      <div className="font-medium">{bead.id}</div>
      <div className="text-sm text-subtext0">{bead.title}</div>
      {bead.assignee && (
        <AgentBadge agent={bead.assignee} size="sm" />
      )}
    </div>
  );
}
```

### 7.3 Triage Panel

```tsx
// components/beads/triage-panel.tsx
export function TriagePanel() {
  const { data: triage, isLoading } = useTriage();

  if (isLoading) return <TriageSkeleton />;

  return (
    <div className="space-y-6">
      {/* Quick Reference */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Reference</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-4 gap-4">
            <Stat label="Open" value={triage.quickRef.openCount} />
            <Stat label="Ready" value={triage.quickRef.actionableCount} />
            <Stat label="Blocked" value={triage.quickRef.blockedCount} />
            <Stat label="In Progress" value={triage.quickRef.inProgressCount} />
          </div>
        </CardContent>
      </Card>

      {/* Top Picks */}
      <Card>
        <CardHeader>
          <CardTitle>Top Picks</CardTitle>
          <CardDescription>Highest impact items to work on next</CardDescription>
        </CardHeader>
        <CardContent>
          {triage.recommendations.slice(0, 5).map(rec => (
            <RecommendationRow
              key={rec.id}
              recommendation={rec}
              onClaim={() => claimBead(rec.id)}
            />
          ))}
        </CardContent>
      </Card>

      {/* Quick Wins */}
      <Card>
        <CardHeader>
          <CardTitle>⚡ Quick Wins</CardTitle>
        </CardHeader>
        <CardContent>
          {triage.quickWins.map(qw => (
            <QuickWinRow key={qw.id} item={qw} />
          ))}
        </CardContent>
      </Card>

      {/* Blockers to Clear */}
      <Card className="border-red-500/50">
        <CardHeader>
          <CardTitle>🚧 Blockers to Clear</CardTitle>
        </CardHeader>
        <CardContent>
          {triage.blockersToClear.map(blocker => (
            <BlockerRow key={blocker.id} blocker={blocker} />
          ))}
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## 8. CASS & Memory System Integration

### 8.1 Memory Search UI

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MEMORY & CONTEXT                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 Search past sessions and memory...                              │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  Filters: [Agent: All ▾] [Workspace: myproject ▾] [Since: 7 days ▾]        │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  RELEVANT MEMORIES                                                   │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  📋 Rule: Always validate JWT expiry before trusting claims          │   │
│  │     Category: authentication | Confidence: 92%                       │   │
│  │     Source: Session 2025-01-05 (auth-feature)                        │   │
│  │     [Apply to Current Task]                                          │   │
│  │                                                                      │   │
│  │  📋 Rule: Use bcrypt with cost factor 12 for password hashing        │   │
│  │     Category: security | Confidence: 88%                             │   │
│  │     Source: Session 2025-01-03 (user-signup)                         │   │
│  │     [Apply to Current Task]                                          │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  PAST SESSIONS                                                       │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  🟣 Claude | auth-feature | 2025-01-05 | Score: 0.94                │   │
│  │     "Implemented JWT validation with refresh token rotation"         │   │
│  │     [View Session] [Extract Rules]                                   │   │
│  │                                                                      │   │
│  │  🔵 Codex | api-endpoints | 2025-01-04 | Score: 0.87                │   │
│  │     "Created REST endpoints for user CRUD operations"                │   │
│  │     [View Session] [Extract Rules]                                   │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  [Refresh Memory]  [Trigger Reflection]  [Export Rules]                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.2 Context Injection Flow

```tsx
// hooks/use-memory-context.ts
export function useMemoryContext(taskDescription: string) {
  return useQuery({
    queryKey: ['memory', 'context', taskDescription],
    queryFn: async () => {
      const [cassContext, cmContext] = await Promise.all([
        api.get('/cass/context', { query: taskDescription }),
        api.get('/memory/context', { task: taskDescription }),
      ]);

      return {
        relevantSessions: cassContext.hits,
        applicableRules: cmContext.rules,
        antiPatterns: cmContext.antiPatterns,
        suggestedQueries: cmContext.suggestedQueries,
      };
    },
  });
}

// Usage in send prompt
function SendPromptWithContext({ session }: { session: string }) {
  const [prompt, setPrompt] = useState('');
  const { data: context } = useMemoryContext(prompt);

  const handleSend = async () => {
    // Inject context into prompt
    const enrichedPrompt = buildEnrichedPrompt(prompt, context);
    await api.post(`/sessions/${session}/send`, {
      prompt: enrichedPrompt,
      include_context: true,
    });
  };

  return (
    <div>
      <Textarea value={prompt} onChange={e => setPrompt(e.target.value)} />

      {/* Show relevant context */}
      {context?.applicableRules.length > 0 && (
        <ContextSuggestions rules={context.applicableRules} />
      )}

      <Button onClick={handleSend}>Send with Context</Button>
    </div>
  );
}
```

---

## 9. UBS Scanner Integration

### 9.1 Scanner Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CODE QUALITY DASHBOARD                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Project: myproject              Last Scan: 2 minutes ago                   │
│  Auto-scan: ● Enabled            Status: ✓ Passing                          │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    SEVERITY BREAKDOWN                                │   │
│  │                                                                      │   │
│  │    Critical  ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0        │   │
│  │    Warning   ████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  3        │   │
│  │    Info      ████████████████████████████████░░░░░░░░░░░░  12       │   │
│  │                                                                      │   │
│  │    Files Scanned: 127        Duration: 3.2s                         │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  FINDINGS                                              [Filter ▾]    │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  ⚠️  src/auth/jwt.go:42                                              │   │
│  │      Unused variable 'tokenExpiry'                                   │   │
│  │      💡 Remove or use the variable                                   │   │
│  │      [Go to File] [Create Bead] [Dismiss]                           │   │
│  │                                                                      │   │
│  │  ⚠️  src/api/users.go:87                                             │   │
│  │      Error return value not checked                                  │   │
│  │      💡 Handle the error or explicitly ignore with _                 │   │
│  │      [Go to File] [Create Bead] [Dismiss]                           │   │
│  │                                                                      │   │
│  │  ℹ️  src/models/user.go:15                                           │   │
│  │      Consider adding JSON tags to struct fields                      │   │
│  │      [Go to File] [Create Bead] [Dismiss]                           │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  [Run Full Scan]  [Scan Staged Only]  [Configure Thresholds]               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Auto-Scanner WebSocket Integration

```tsx
// hooks/use-scanner.ts
export function useAutoScanner(projectDir: string) {
  const queryClient = useQueryClient();

  // Subscribe to scanner channel
  useWebSocketChannel('scanner', (event) => {
    if (event.payload.event === 'scan_complete') {
      // Update cache with new results
      queryClient.setQueryData(
        ['scanner', 'latest', projectDir],
        event.payload
      );

      // Show notification for critical findings
      if (event.payload.totals.critical > 0) {
        toast.error(`${event.payload.totals.critical} critical issues found!`);
      }
    }
  });

  return useQuery({
    queryKey: ['scanner', 'latest', projectDir],
    queryFn: () => api.get('/scanner/latest', { project: projectDir }),
  });
}

// Auto-create beads from findings
async function createBeadFromFinding(finding: Finding) {
  await api.post('/beads', {
    title: `Fix: ${finding.message}`,
    type: 'bug',
    priority: finding.severity === 'critical' ? 0 : 2,
    labels: ['ubs', finding.category],
    description: `
**File:** ${finding.file}:${finding.line}
**Severity:** ${finding.severity}
**Message:** ${finding.message}
**Suggestion:** ${finding.suggestion}
    `,
  });
}
```

---

## 10. CAAM Account Management

### 10.1 Account Manager UI

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ACCOUNT MANAGER                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Auto-Rotate: ● Enabled          Reset Buffer: 15 minutes                   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  CLAUDE ACCOUNTS                                                     │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  ✓ primary (active)                                                  │   │
│  │    user@example.com | Priority: 1 | Rate Limit: OK                   │   │
│  │    [Deactivate]                                                      │   │
│  │                                                                      │   │
│  │  ○ backup-1                                                          │   │
│  │    backup@example.com | Priority: 2 | Cooldown: 12m                  │   │
│  │    [Activate]                                                        │   │
│  │                                                                      │   │
│  │  ○ backup-2                                                          │   │
│  │    team@example.com | Priority: 3 | Rate Limit: OK                   │   │
│  │    [Activate]                                                        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  CODEX ACCOUNTS                                                      │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │                                                                      │   │
│  │  ✓ main (active)                                                     │   │
│  │    org-xxx | Priority: 1 | Rate Limit: WARNING (80%)                │   │
│  │    [Rotate Now]                                                      │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  GEMINI ACCOUNTS                                                     │   │
│  ├─────────────────────────────────────────────────────────────────────┤   │
│  │  ✓ default (active)                                                  │   │
│  │    project-xxx | Rate Limit: OK                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  [Add Account]  [Import from CAAM]  [Configure Auto-Rotate]                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 10.2 Rate Limit Handling

```tsx
// hooks/use-accounts.ts
export function useAccountRotation() {
  const queryClient = useQueryClient();

  // Listen for rate limit events
  useWebSocketChannel('accounts', (event) => {
    if (event.payload.event === 'rate_limit_hit') {
      // Show notification
      toast.warning(
        `${event.payload.agent_type} hit rate limit. Rotating to ${event.payload.new_account}...`
      );

      // Invalidate account queries
      queryClient.invalidateQueries({ queryKey: ['accounts'] });
    }

    if (event.payload.event === 'account_rotated') {
      toast.success(
        `Rotated to ${event.payload.new_account}. Cooldown until ${formatTime(event.payload.cooldown_until)}`
      );
    }
  });

  return useMutation({
    mutationFn: (params: { type: string; account: string }) =>
      api.post(`/accounts/${params.type}/activate`, { account: params.account }),
  });
}
```

---

## 11. SLB Safety Guardrails

### 11.1 Two-Person Rule UI

When dangerous commands require approval:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      ⚠️ APPROVAL REQUIRED                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Agent "GreenCastle" is requesting approval for a dangerous operation:      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  Command: rm -rf ./node_modules && npm install                       │   │
│  │  Risk Level: 🟡 Medium                                               │   │
│  │  Working Directory: /data/projects/myproject                         │   │
│  │  Requested: 2 minutes ago                                            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  Context:                                                                   │
│  "Attempting to fix corrupted node_modules after merge conflict"            │
│                                                                             │
│  Similar past commands:                                                     │
│  • Approved 3 times in last 7 days                                         │
│  • Never caused issues                                                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  [Approve]  [Approve All Similar]  [Deny]  [Deny with Message]      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ☐ Remember this decision for similar commands                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 11.2 Safety Dashboard

```tsx
// components/safety/safety-dashboard.tsx
export function SafetyDashboard() {
  const { data: pendingApprovals } = usePendingApprovals();
  const { data: recentDecisions } = useRecentDecisions();
  const { data: safetyConfig } = useSafetyConfig();

  return (
    <div className="space-y-6">
      {/* Mode Toggle */}
      <Card>
        <CardHeader>
          <CardTitle>Safety Mode</CardTitle>
        </CardHeader>
        <CardContent>
          <RadioGroup value={safetyConfig.mode}>
            <RadioGroupItem value="vibe">
              <Label>Vibe Mode (Throwaway VPS)</Label>
              <Description>Agents have dangerous permissions</Description>
            </RadioGroupItem>
            <RadioGroupItem value="safe">
              <Label>Safe Mode (Production-like)</Label>
              <Description>Two-person rule for dangerous commands</Description>
            </RadioGroupItem>
          </RadioGroup>
        </CardContent>
      </Card>

      {/* Pending Approvals */}
      {pendingApprovals?.length > 0 && (
        <Card className="border-yellow-500">
          <CardHeader>
            <CardTitle>Pending Approvals ({pendingApprovals.length})</CardTitle>
          </CardHeader>
          <CardContent>
            {pendingApprovals.map(approval => (
              <ApprovalRequest
                key={approval.id}
                approval={approval}
                onApprove={() => approve(approval.id)}
                onDeny={() => deny(approval.id)}
              />
            ))}
          </CardContent>
        </Card>
      )}

      {/* Recent Decisions */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Decisions</CardTitle>
        </CardHeader>
        <CardContent>
          <DecisionHistory decisions={recentDecisions} />
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## 12. Pipeline & Workflow Engine

### 12.1 Visual Pipeline Builder

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PIPELINE BUILDER                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Pipeline: auth-feature-workflow                        [Save] [Run] [...]  │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                                                                      │   │
│  │    ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐  │   │
│  │    │  Start  │─────▶│  Spawn  │─────▶│  Send   │─────▶│  Wait   │  │   │
│  │    │         │      │ Agents  │      │ Prompt  │      │  Idle   │  │   │
│  │    └─────────┘      └─────────┘      └─────────┘      └────┬────┘  │   │
│  │                                                             │       │   │
│  │                          ┌──────────────────────────────────┘       │   │
│  │                          ▼                                          │   │
│  │    ┌─────────┐      ┌─────────┐      ┌─────────┐      ┌─────────┐  │   │
│  │    │   End   │◀─────│ Deploy  │◀─────│  Test   │◀─────│  Scan   │  │   │
│  │    │         │      │         │      │         │      │  (UBS)  │  │   │
│  │    └─────────┘      └─────────┘      └─────────┘      └─────────┘  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  STEP PALETTE                                                        │   │
│  │                                                                      │   │
│  │  [Spawn Agents] [Send Prompt] [Wait] [Scan] [Checkpoint]            │   │
│  │  [Condition] [Loop] [Parallel] [Call Pipeline] [Shell]              │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │  STEP EDITOR: Spawn Agents                                          │   │
│  │                                                                      │   │
│  │  Session: ${session_name}                                           │   │
│  │  Claude Count: [3]                                                   │   │
│  │  Codex Count: [2]                                                   │   │
│  │  Gemini Count: [1]                                                  │   │
│  │  Recipe: [full-stack ▾]                                             │   │
│  │  Auto-restart: [✓]                                                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 12.2 Pipeline Execution Monitor

```tsx
// components/pipeline/execution-monitor.tsx
export function PipelineExecutionMonitor({ pipelineId }: { pipelineId: string }) {
  const { data: pipeline } = usePipeline(pipelineId);
  const { data: steps } = usePipelineSteps(pipelineId);

  // Real-time updates via WebSocket
  useWebSocketChannel('pipeline', (event) => {
    if (event.payload.pipeline_id === pipelineId) {
      queryClient.setQueryData(['pipeline', pipelineId], (old) => ({
        ...old,
        ...event.payload,
      }));
    }
  });

  return (
    <div className="space-y-4">
      {/* Progress Bar */}
      <div className="relative">
        <Progress value={pipeline.progress.percent} />
        <span className="absolute right-0 top-0 text-sm">
          {pipeline.progress.completed}/{pipeline.progress.total}
        </span>
      </div>

      {/* Step List */}
      <div className="space-y-2">
        {steps.map((step, index) => (
          <StepCard
            key={step.id}
            step={step}
            isActive={index === pipeline.currentStep}
            isComplete={step.status === 'success'}
            isFailed={step.status === 'failed'}
          />
        ))}
      </div>

      {/* Actions */}
      <div className="flex gap-2">
        {pipeline.status === 'running' && (
          <Button variant="destructive" onClick={() => cancelPipeline(pipelineId)}>
            Cancel
          </Button>
        )}
        {pipeline.status === 'failed' && (
          <Button onClick={() => retryPipeline(pipelineId)}>
            Retry from Failed Step
          </Button>
        )}
      </div>
    </div>
  );
}
```

---

## 13. Web UI Layer

### 13.1 Technology Stack

| Category | Technology | Version |
|----------|------------|---------|
| **Framework** | Next.js | 16.x |
| **React** | React | 19.2 |
| **Language** | TypeScript | 5.5+ (strict mode) |
| **Styling** | Tailwind CSS | 4.x |
| **Animation** | Framer Motion | 11.x |
| **Icons** | Lucide React | Latest |
| **State (Server)** | TanStack Query | 5.x |
| **State (Client)** | Zustand | 5.x |
| **Forms** | TanStack Form | 1.x |
| **Tables** | TanStack Table | 8.x |
| **Charts** | Recharts | 2.x |
| **Terminal** | xterm.js | 5.x |
| **Graphs** | React Flow | 11.x |
| **Package Manager** | Bun | 1.x |

### 13.2 Design System

#### Color Palette (Catppuccin Mocha)

```typescript
export const catppuccinMocha = {
  // Base colors
  base: '#1e1e2e',
  mantle: '#181825',
  crust: '#11111b',

  // Surface colors
  surface0: '#313244',
  surface1: '#45475a',
  surface2: '#585b70',

  // Text colors
  text: '#cdd6f4',
  subtext1: '#bac2de',
  subtext0: '#a6adc8',

  // Agent colors (matching NTM TUI)
  mauve: '#cba6f7',    // Claude (purple)
  blue: '#89b4fa',     // Codex (blue)
  yellow: '#f9e2af',   // Gemini (yellow)
  green: '#a6e3a1',    // User/Success
  red: '#f38ba8',      // Error/Critical
  peach: '#fab387',    // Warning
  teal: '#94e2d5',     // Info
};
```

#### Animation Principles

```typescript
export const fadeIn = {
  initial: { opacity: 0 },
  animate: { opacity: 1 },
  exit: { opacity: 0 },
  transition: { duration: 0.2 }
};

export const slideUp = {
  initial: { opacity: 0, y: 20 },
  animate: { opacity: 1, y: 0 },
  exit: { opacity: 0, y: -20 },
  transition: { duration: 0.3, ease: [0.4, 0, 0.2, 1] }
};

export const staggerChildren = {
  animate: { transition: { staggerChildren: 0.05 } }
};
```

### 13.3 Page Structure

```
app/
├── (auth)/
│   ├── dashboard/           # Main dashboard
│   │   └── page.tsx
│   ├── sessions/
│   │   ├── page.tsx         # Session list
│   │   └── [name]/
│   │       ├── page.tsx     # Session detail
│   │       ├── terminal/    # Full terminal view
│   │       └── output/      # Output viewer
│   ├── beads/
│   │   ├── page.tsx         # Kanban board
│   │   ├── graph/           # Dependency graph
│   │   ├── triage/          # Triage panel
│   │   └── [id]/            # Bead detail
│   ├── mail/
│   │   ├── page.tsx         # Unified inbox
│   │   ├── threads/[id]/    # Thread view
│   │   └── reservations/    # File reservation map
│   ├── memory/
│   │   ├── page.tsx         # Search UI
│   │   ├── playbook/        # Rule browser
│   │   └── sessions/        # Session viewer
│   ├── scanner/
│   │   └── page.tsx         # Scanner dashboard
│   ├── accounts/
│   │   └── page.tsx         # Account manager
│   ├── pipelines/
│   │   ├── page.tsx         # Pipeline list
│   │   ├── builder/         # Visual builder
│   │   └── [id]/            # Execution monitor
│   ├── safety/
│   │   └── page.tsx         # Safety dashboard
│   ├── palette/
│   │   └── page.tsx         # Command palette
│   └── settings/
│       └── page.tsx         # Configuration
├── layout.tsx               # Root layout
└── page.tsx                 # Landing/login
```

---

## 14. Desktop vs Mobile UX Strategy

### 14.1 Philosophy

Desktop and mobile are fundamentally different interaction paradigms. We design **separate optimized experiences**.

### 14.2 Desktop Experience (≥1024px)

```
┌────────────────────────────────────────────────────────────────────────────┐
│  [Logo]  Session: myproject  │  [Search ⌘K]     [Alerts 3] [Settings]     │
├──────────────┬─────────────────────────────────────────────────────────────┤
│              │                                                             │
│  Sessions    │   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐      │
│  ─────────   │   │ cc_1    │  │ cc_2    │  │ cc_3    │  │ cod_1   │      │
│  • myproject │   │ Claude  │  │ Claude  │  │ Claude  │  │ Codex   │      │
│  • backend   │   │ Green   │  │ Blue    │  │ Red     │  │         │      │
│              │   │ Castle  │  │ Lake    │  │ Stone   │  │         │      │
│  Beads       │   │ ████░░  │  │ ██████  │  │ ███░░░  │  │ █████░  │      │
│  ─────────   │   └─────────┘  └─────────┘  └─────────┘  └─────────┘      │
│  Ready: 9    │                                                             │
│  Blocked: 3  │   ┌──────────────────────────────────────────────────────┐ │
│              │   │  Terminal: myproject__cc_1 (GreenCastle)              │ │
│  Mail        │   │  ─────────────────────────────────────────────────── │ │
│  ─────────   │   │  $ claude --dangerously-skip-permissions              │ │
│  Unread: 5   │   │  > Analyzing authentication module...                 │ │
│              │   │  > Found 3 potential issues:                          │ │
│  Scanner     │   │  > 1. JWT validation missing expiry check             │ │
│  ─────────   │   │  >                                                    │ │
│  ✓ Passing   │   └──────────────────────────────────────────────────────┘ │
│              │                                                             │
│  Memory      │   ┌──────────────────────────────────────────────────────┐ │
│  ─────────   │   │  File Reservations                                    │ │
│  3 rules     │   │  src/auth/** → GreenCastle (45m)                     │ │
│  available   │   │  src/api/** → BlueLake (30m)                         │ │
│              │   └──────────────────────────────────────────────────────┘ │
├──────────────┴─────────────────────────────────────────────────────────────┤
│  [Target: All ▾]  Enter prompt...                        [+ Context] [Send]│
└────────────────────────────────────────────────────────────────────────────┘
```

**Key Desktop Features:**
- Persistent sidebar with all flywheel tools
- Multi-pane agent view with terminals
- File reservation map always visible
- Command palette (`⌘K`)
- Keyboard navigation
- Drag-and-drop task management

### 14.3 Mobile Experience (<768px)

```
┌────────────────────────────┐
│  ← myproject     [•••]    │
├────────────────────────────┤
│                            │
│  ╭────────────────────╮   │
│  │ 🟣 GreenCastle     │   │
│  │ Claude • Working   │   │
│  │ ████████░░░░ 72%   │   │
│  │                    │   │
│  │ "Analyzing auth    │   │
│  │  module..."        │   │
│  ╰────────────────────╯   │
│                            │
│  ╭────────────────────╮   │
│  │ 🔵 BlueLake        │   │
│  │ Codex • Idle       │   │
│  │ █████░░░░░░░ 45%   │   │
│  ╰────────────────────╯   │
│                            │
│  ╭────────────────────╮   │
│  │ 🟡 RedStone        │   │
│  │ Gemini • Working   │   │
│  │ ██░░░░░░░░░░ 23%   │   │
│  ╰────────────────────╯   │
│                            │
├────────────────────────────┤
│  ╭────────────────────╮   │
│  │ Quick message...   │ ↑ │
│  ╰────────────────────╯   │
├────────────────────────────┤
│ 🏠  📋  ✉️  🔍  ⚙️        │
│Home Beads Mail Memory Set │
└────────────────────────────┘
```

**Key Mobile Features:**
- Bottom navigation for flywheel tools
- Swipe gestures on agent cards
- Pull-to-refresh
- Simplified terminal (read-only)
- Quick prompts instead of keyboard
- Push notifications for alerts
- Haptic feedback

---

## 15. Agent SDK Integration Strategy

### 15.1 Dual-Mode Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     AGENT EXECUTION MODES                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────┐      ┌─────────────────────────┐      │
│  │    TMUX MODE        │      │      SDK MODE           │      │
│  │    (Default)        │      │    (Alternative)        │      │
│  ├─────────────────────┤      ├─────────────────────────┤      │
│  │                     │      │                         │      │
│  │  • Visual terminal  │      │  • Programmatic API     │      │
│  │  • Full scrollback  │      │  • Structured events    │      │
│  │  • Multi-session    │      │  • Lower overhead       │      │
│  │  • SSH support      │      │  • No tmux dependency   │      │
│  │                     │      │                         │      │
│  │  Uses:              │      │  Uses:                  │      │
│  │  • tmux panes       │      │  • claude-agent-sdk     │      │
│  │  • Text streams     │      │  • codex-sdk            │      │
│  │  • CLI spawning     │      │  • @google/genai        │      │
│  └─────────────────────┘      └─────────────────────────┘      │
│                                                                 │
│                    ▼ Unified Event Interface ▼                  │
│                                                                 │
│  Both modes emit: output, tool_call, tool_result,              │
│                   status, error, complete                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 15.2 ACP Integration (Future)

When Agent Client Protocol becomes more widely adopted:

```typescript
// Future ACP runner
class ACPRunner implements AgentRunner {
  private client: ACPClient;

  async *run(prompt: string): AsyncGenerator<AgentEvent> {
    const requestId = this.client.sendRequest('agent/run', { prompt });
    for await (const message of this.client.responses(requestId)) {
      yield this.transformACPMessage(message);
    }
  }
}
```

---

## 16. Implementation Phases

### Phase 1: Foundation (Weeks 1-3)

**Goal**: Core REST API and WebSocket infrastructure

- [ ] Create `internal/api` package with Chi router
- [ ] Implement session CRUD endpoints
- [ ] Implement WebSocket hub with output channel
- [ ] Add `ntm serve` command
- [ ] Generate OpenAPI spec
- [ ] Add Swagger UI

### Phase 2: Flywheel Tool APIs (Weeks 4-6)

**Goal**: Full API coverage for all 8 tools

- [ ] Agent Mail endpoints (messaging, reservations, contacts)
- [ ] Beads endpoints (CRUD, triage, insights)
- [ ] CASS endpoints (search, context)
- [ ] Scanner endpoints (scan, auto-scan)
- [ ] Account endpoints (list, activate)
- [ ] Pipeline endpoints (run, status, cancel)
- [ ] All WebSocket channels

### Phase 3: Web UI Foundation (Weeks 7-10)

**Goal**: Core Next.js application

- [ ] Initialize Next.js 16 project with Bun
- [ ] Implement design system (Catppuccin theme)
- [ ] Build API client with TanStack Query
- [ ] Create WebSocket provider
- [ ] Build dashboard page
- [ ] Build session detail with terminal

### Phase 4: Flywheel Tool UIs (Weeks 11-16)

**Goal**: UI for each flywheel tool

- [ ] Beads board (Kanban, graph, triage)
- [ ] Agent Mail UI (inbox, threads, reservations)
- [ ] Memory search UI
- [ ] Scanner dashboard
- [ ] Account manager
- [ ] Pipeline builder
- [ ] Safety dashboard

### Phase 5: Mobile & Polish (Weeks 17-20)

**Goal**: Mobile optimization and production readiness

- [ ] Mobile-specific layouts
- [ ] Touch interactions
- [ ] Push notifications
- [ ] Performance optimization
- [ ] Accessibility audit
- [ ] Documentation

### Phase 6: SDK Mode (Weeks 21-24)

**Goal**: Alternative execution mode

- [ ] Claude Agent SDK runner
- [ ] Codex SDK runner
- [ ] Gemini SDK runner
- [ ] Mode selection UI
- [ ] ACP exploration

---

## 17. File Structure

### Go Backend Additions

```
internal/
├── api/                          # REST API layer
│   ├── api.go
│   ├── handlers/
│   │   ├── sessions.go
│   │   ├── beads.go
│   │   ├── mail.go
│   │   ├── reservations.go
│   │   ├── cass.go
│   │   ├── memory.go
│   │   ├── scanner.go
│   │   ├── accounts.go
│   │   ├── pipelines.go
│   │   ├── safety.go
│   │   └── system.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logging.go
│   └── openapi/
├── ws/                           # WebSocket layer
│   ├── hub.go
│   ├── client.go
│   ├── channels/
│   │   ├── output.go
│   │   ├── beads.go
│   │   ├── mail.go
│   │   ├── scanner.go
│   │   ├── pipeline.go
│   │   └── ...
│   └── protocol.go
└── serve/
    └── server.go
```

### Frontend Structure

```
web/
├── app/
│   ├── (auth)/
│   │   ├── dashboard/
│   │   ├── sessions/
│   │   ├── beads/
│   │   ├── mail/
│   │   ├── memory/
│   │   ├── scanner/
│   │   ├── accounts/
│   │   ├── pipelines/
│   │   ├── safety/
│   │   └── settings/
│   └── layout.tsx
├── components/
│   ├── ui/
│   ├── dashboard/
│   ├── session/
│   ├── beads/
│   ├── mail/
│   ├── memory/
│   ├── scanner/
│   ├── pipeline/
│   ├── terminal/
│   └── mobile/
├── hooks/
├── lib/
├── stores/
└── types/
```

---

## 18. Technical Specifications

### API Performance Targets

| Metric | Target |
|--------|--------|
| REST response time (p50) | < 50ms |
| REST response time (p99) | < 200ms |
| WebSocket latency | < 100ms |
| Concurrent WebSocket connections | 1000+ |

### Frontend Performance Targets

| Metric | Target |
|--------|--------|
| First Contentful Paint | < 1.2s |
| Largest Contentful Paint | < 2.5s |
| Time to Interactive | < 3.5s |
| Cumulative Layout Shift | < 0.1 |

### Browser Support

| Browser | Minimum Version |
|---------|-----------------|
| Chrome | 111+ |
| Firefox | 111+ |
| Safari | 16.4+ |
| Edge | 111+ |

---

## Appendix A: Research Sources

- [Agent Client Protocol](https://agentclientprotocol.com/)
- [Claude Code in Zed via ACP](https://zed.dev/blog/claude-code-via-acp)
- [@anthropic-ai/claude-agent-sdk](https://www.npmjs.com/package/@anthropic-ai/claude-agent-sdk)
- [@openai/codex-sdk](https://developers.openai.com/codex/sdk/)
- [@google/genai](https://www.npmjs.com/package/@google/genai)
- [Next.js 16](https://nextjs.org/blog/next-16)
- [TanStack Query + WebSockets](https://tkdodo.eu/blog/using-web-sockets-with-react-query)
- [WebSocket Architecture Best Practices](https://ably.com/topic/websocket-architecture-best-practices)
- [Stripe Apps UI Toolkit](https://docs.stripe.com/stripe-apps/components)

---

*Document Version: 2.0.0*
*Last Updated: January 7, 2025*
*Author: Claude Opus 4.5*
