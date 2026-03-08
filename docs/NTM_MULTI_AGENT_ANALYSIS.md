# NTM 作为多代理主控工具分析

> 版本: 1.0
> 日期: 2026-01-27
> 分析范围: NTM v1.5.0+ (commit ca1d752)

---

## 一、执行摘要

**结论**: NTM **适合**做多代理同时工作的主控工具，但有一些重要限制需要注意。

| 维度 | 评分 | 说明 |
|------|------|------|
| 代理管理 | ⭐⭐⭐⭐⭐ | 优秀的 tmux 集成，支持 Claude/Codex/Gemini |
| 任务分发 | ⭐⭐⭐⭐ | 广播发送，智能路由已实现 |
| 状态感知 | ⭐⭐⭐⭐ | Activity Detection、Health Monitoring |
| 工作流编排 | ⭐⭐⭐⭐⭐ | Pipeline 支持 DAG、并行、循环、条件 |
| 通信协调 | ⭐⭐⭐⭐ | Agent Mail MCP 集成，消息系统 |
| 冲突避免 | ⭐⭐⭐⭐ | File reservations + staggered spawn |

---

## 二、核心能力

### 1. 代理生命周期管理

```bash
# 创建多代理会话
ntm spawn myproject --cc=2 --cod=1 --gmi=1

# 添加更多代理
ntm add myproject --cc=1

# 监控状态
ntm activity myproject
ntm health myproject
```

**支持**:
- Claude Code (cc)
- Codex (cod)
- Gemini (gmi)
- 任意 CLI 工具

### 2. 任务分发策略

| 策略 | 命令 | 适用场景 |
|------|------|----------|
| 广播 | `ntm send --all` | 需要多视角分析 |
| 类型过滤 | `ntm send --cc` | 特定能力代理 |
| 智能路由 | `ntm send --smart` | 自动选择最优代理 |
| 指定 pane | `ntm send --pane=1` | 精确控制 |

**智能路由评分算法**:
```
Score = (context_score × 0.4) + (state_score × 0.4) + (recency_score × 0.2)
```

### 3. Pipeline 工作流

```yaml
id: parallel_workflow
steps:
  - id: research
    parallel:
      - id: market
        agent: claude
        prompt: Research market
      - id: tech
        agent: codex
        prompt: Research tech
      - id: competitor
        agent: gemini
        prompt: Analyze competitors

  - id: synthesis
    depends_on: [research]
    prompt: Combine findings
```

**特性**:
- DAG 依赖图
- 并行执行 (最多 8 个并行)
- 循环 (for-each, while, times)
- 条件执行 (when)
- 输出解析 (JSON, YAML, regex)
- 错误处理 (fail, continue, retry)

### 4. 协调机制

#### Agent Mail (MCP 集成)
- 代理身份注册
- Inbox/Outbox 消息
- File reservations (避免冲突)
- Thread 协作

#### Staggered Spawn (防惊群)
```bash
ntm spawn --cc=3 --stagger  # 90s 间隔
```

解决多代理同时抢夺同一任务的竞态问题。

---

## 三、优势

### 1. 基础设施成熟

- **Tmux 深度集成**: 基于 tmux pane 管理，稳定可靠
- **Robot Mode API**: 完整的 JSON API，可编程控制
- **跨平台**: Go 编写，单二进制，无依赖

### 2. 状态感知能力强

```
┌──────────────┬────────────────────────────────────────────────────────┐
│ State        │ Meaning                                                │
├──────────────┼────────────────────────────────────────────────────────┤
│ GENERATING   │ 高速输出 (>10 c/s)                                      │
│ WAITING      │ 空闲提示，可接收任务                                     │
│ THINKING     │ 低速处理，思考中                                        │
│ ERROR        │ 错误检测                                               │
│ STALLED      │ 无活动，可能卡死                                        │
└──────────────┴────────────────────────────────────────────────────────┘
```

### 3. 工作流能力完整

- 依赖拓扑排序 (Kahn 算法)
- 循环依赖检测
- 状态持久化 (`.ntm/pipelines/*.json`)
- 失败恢复

### 4. 生态集成

| 工具 | 集成方式 | 用途 |
|------|----------|------|
| CASS | cli | 历史会话搜索 |
| bv (beads-visualizer) | cli | 任务优先级/依赖分析 |
| br (beads_rust) | cli | Issue 追踪 |
| Agent Mail | MCP | 代理间通信 |

---

## 四、劣势与限制

### 1. 架构限制

| 问题 | 说明 | 影响 |
|------|------|------|
| **本地集中式** | 必须在一台机器上运行 | 无法跨机器分布式 |
| **Tmux 依赖** | 需要 tmux 环境 | Windows 支持弱 |
| **同步交互** | 基于轮询，非事件驱动 | 延迟较高 |

### 2. 已知问题 (来自你提的 Issues)

#### Issue #20: Pane 标题冲突 (已修复)
- **问题**: Claude Code 动态修改 pane 标题，NTM 无法识别
- **修复**: 添加进程名 fallback 检测

#### Issue #23: Agent Health Detection (已修复)
- **问题**: 健康检测不准确，bulk-assign 失败
- **修复**: 改进检测逻辑

#### Issue #24: Agent Identity Persistence (已修复)
- **问题**: MCP response envelope 解析错误，agent_name 为空
- **修复**: 正确提取 `structuredContent`

#### Issue #24 遗留问题: AGENT_NAME Fragility
- **问题**: 即使设置 `AGENT_NAME` env var，cc/cod 调用 `macro_start_session` 不会自动读取
- **Workaround**: 手动 send 告知身份
- **状态**: 未解决，破坏自动化

### 3. 复杂度

```bash
# 命令数量
ntm --help | grep "Available Commands" -A 100 | wc -l
# > 50 个子命令
```

学习曲线陡峭，需要理解:
- tmux 基础
- ntm 命令体系
- Agent Mail MCP 协议
- beads/bv/br 工具链

### 4. 缺少的功能

| 功能 | 状态 | 替代方案 |
|------|------|----------|
| Web UI | 规划中 | CLI/TUI |
| 分布式 | 不支持 | 单机多会话 |
| 动态扩缩容 | 部分支持 | `ntm add` 手动 |
| 结果聚合 | 手动 | `ntm diff` 对比 |

---

## 五、适用场景

### ✅ 推荐场景

1. **单机多代理并行开发**
   - 同时代码审查、测试、文档编写
   - 多角度问题分析

2. **复杂工作流自动化**
   - 设计 → 实现 → 测试 → 审查 pipeline
   - 循环处理批量任务

3. **研究/实验环境**
   - 需要观察每个代理输出
   - 手动介入调试

### ⚠️ 谨慎使用

1. **生产环境自动化**
   - AGENT_NAME fragility 问题
   - 需要额外监控

2. **大规模部署**
   - 单机限制
   - 手动扩容

### ❌ 不推荐

1. **跨机器分布式任务**
   - 架构不支持

2. **简单任务**
   - 复杂度不值得

---

## 六、替代方案对比

| 工具 | 架构 | 优点 | 缺点 |
|------|------|------|------|
| **NTM** | 本地 tmux | 可视化、完整工具链 | 单机、复杂 |
| **LangGraph** | 分布式 DAG | 云原生、可扩展 | 需要 infra |
| **AutoGen** | Python 库 | 灵活、轻量 | 无 UI、需编程 |
| **Chorus** (NTM子项目) | Agent Mail 消息 | 松耦合 | 实验性 |

---

## 七、建议

### 如果选择 NTM

1. **初始化**
   ```bash
   ntm init
   eval "$(ntm shell zsh)"  # 或 bash
   ```

2. **使用 Staggered Spawn**
   ```bash
   ntm spawn myproject --cc=3 --stagger
   ```

3. **利用 Pipeline**
   ```yaml
   # .ntm/workflows/standard.yaml
   steps:
     - id: triage
       agent: claude
       prompt: "bv --robot-triage"
   ```

4. **监控健康**
   ```bash
   ntm health myproject --watch
   ```

### 改进方向

1. **高优先级**
   - 修复 AGENT_NAME 自动注入
   - Web UI (简化使用)
   - 更好的错误恢复

2. **中优先级**
   - 分布式支持 (Agent Mail 远程)
   - 自动扩缩容
   - 结果聚合 AI

---

## 八、结论

**NTM 是目前最成熟的开源多代理主控工具之一**，特别适合:

- 需要可视化监控的场景
- 单机多代理协作
- 复杂工作流编排

但需要权衡:
- 学习成本
- 单机限制
- 某些自动化问题需 workaround

**推荐用于**: 研发团队的高级用户、工具链完善的项目

**不推荐用于**: 简单任务、需要分布式部署的场景
