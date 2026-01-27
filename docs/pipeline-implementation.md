# NTM Pipeline 实现原理

## 概述

NTM Pipeline 是一个为 AI 代理编排设计的复杂工作流系统,支持多步骤、多代理的协作任务执行。它采用依赖驱动的执行模型,具备状态持久化、错误恢复和并行执行能力。

## 核心概念

### 1. 设计原则

- **基于代理的执行**: 每个步骤可以指定特定的 AI 代理(Claude、Codex、Gemini),或让系统根据能力自动路由
- **依赖驱动**: 步骤基于依赖图执行,支持在可能的情况下并行执行
- **状态持久化**: 每个步骤执行后检查点持久化状态,支持失败后恢复
- **变量替换**: 动态变量解析允许步骤间共享数据和上下文
- **错误处理**: 可配置的错误策略(fail、continue、retry、fail-fast)

### 2. 架构组件

```
┌─────────────────────────────────────────────────────────────┐
│                    Workflow Executor                        │
│  (internal/pipeline/executor.go - 1,751 lines)              │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Dependency  │  │  Loop        │  │  Variable    │     │
│  │  Graph       │  │  Executor    │  │  Resolver    │     │
│  │  (deps.go)   │  │  (loops.go)  │  │  (vars.go)   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
├─────────────────────────────────────────────────────────────┤
│                    Agent Routing                            │
│  (internal/robot/routing.go)                               │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Scoring    │  │   Health     │  │  Scheduling  │     │
│  │   Engine     │  │   Checks     │  │  Strategy    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## 核心组件详解

### 1. Workflow Executor (执行引擎)

**文件位置**: `internal/pipeline/executor.go`

这是 pipeline 的核心执行引擎,负责:

- **工作流加载和验证**: 解析 YAML/TOML 格式的工作流文件
- **依赖解析**: 构建和验证依赖图,检测循环依赖
- **执行调度**: 根据依赖关系调度步骤执行
- **状态管理**: 持久化执行状态,支持恢复
- **错误处理**: 实现重试逻辑和错误恢复
- **并行执行**: 支持并发步骤执行(默认最多 8 个并行)

**关键数据结构**:

```go
type Workflow struct {
    ID          string            `yaml:"id"`
    Name        string            `yaml:"name"`
    Description string            `yaml:"description"`
    Steps       []Step            `yaml:"steps"`
    Vars        map[string]any    `yaml:"vars"`
    OnError     string            `yaml:"on_error"` // fail, continue, fail_fast
    MaxParallel int               `yaml:"max_parallel"`
}

type Step struct {
    ID          string                 `yaml:"id"`
    Name        string                 `yaml:"name"`
    Description string                 `yaml:"description"`
    Agent       string                 `yaml:"agent"` // claude, codex, gemini
    Prompt      string                 `yaml:"prompt"`
    DependsOn   []string               `yaml:"depends_on"`
    OnError     string                 `yaml:"on_error"`
    Output      OutputConfig           `yaml:"output"`
    Loop        *LoopConfig            `yaml:"loop"`
    Parallel    *ParallelConfig        `yaml:"parallel"`
    Timeout     string                 `yaml:"timeout"`
    Retry       *RetryConfig           `yaml:"retry"`
}
```

### 2. Dependency Graph (依赖图)

**文件位置**: `internal/pipeline/deps.go`

负责管理步骤间的依赖关系:

- **拓扑排序**: 使用 Kahn 算法计算执行顺序
- **循环检测**: 使用 DFS 检测依赖图中的环
- **并行级别识别**: 确定可以同时执行的步骤
- **失败依赖跟踪**: 当依赖步骤失败时正确跳过下游步骤

**核心算法**:

```go
func (g *Graph) TopologicalSort() ([]string, error) {
    // 计算入度
    inDegree := make(map[string]int)
    for node := range g.Dependencies {
        inDegree[node] = 0
    }
    for _, deps := range g.Dependencies {
        for _, dep := range deps {
            inDegree[dep]++
        }
    }

    // Kahn 算法
    queue := []string{}
    for node, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, node)
        }
    }

    result := []string{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)

        for _, neighbor := range g.Dependencies[node] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != len(g.Dependencies) {
        return nil, errors.New("dependency cycle detected")
    }

    return result, nil
}
```

### 3. 变量替换系统

**文件位置**: `internal/pipeline/vars.go`

支持复杂的变量引用:

- **工作流变量**: `${vars.variable_name}`
- **步骤输出**: `${steps.step_id.output}`
- **嵌套字段**: `${vars.data.nested.field}`
- **默认值**: `${vars.x | "default_value"}`
- **循环变量**: `${loop.item}`, `${loop.index}`

**替换流程**:

1. 解析变量引用模式 `${...}`
2. 根据前缀确定变量类型(vars, steps, loop)
3. 递归解析嵌套字段访问
4. 应用默认值(如果指定)
5. 转换为字符串并替换

### 4. Agent 路由与选择

**文件位置**: `internal/robot/routing.go`

实现智能的代理选择和负载均衡:

#### 评分算法

代理评分基于多个因素:

```go
func calculateScore(agent *AgentState, preference *AgentPreference) float64 {
    score := 0.0

    // 1. 上下文使用 (40%)
    // 分数与使用成反比: score = (1 - usage/maxUsage) * 0.4
    usageScore := (1.0 - float64(agent.ContextUsage)/float64(agent.MaxContext)) * 0.4
    score += usageScore

    // 2. 状态评分 (40%)
    switch agent.State {
    case StateIdle:
        score += 0.4
    case StateWorking:
        score += 0.2
    case StateError:
        score += 0.0  // 错误状态基本不参与选择
    }

    // 3. 最近使用时间 (20%)
    // 越久未使用分数越高
    idleTime := time.Since(agent.LastActivity)
    recencyScore := math.Min(idleTime.Minutes()/10.0, 1.0) * 0.2
    score += recencyScore

    // 4. 亲和性加成
    if preference != nil && preference.PrioritizeAgent == agent.Name {
        score += preference.AffinityBonus
    }

    // 5. 健康检查
    if agent.IsRateLimited || agent.ErrorCount > 5 {
        score *= 0.1  // 严重降低不健康代理的分数
    }

    return score
}
```

#### 路由策略

1. **最少负载 (Least-loaded)**: 选择上下文使用率最低的代理
2. **首选可用 (First-available)**: 快速选择第一个可用代理
3. **轮询 (Round-robin)**: 平均分配请求

### 5. 循环执行器

**文件位置**: `internal/pipeline/loops.go`

支持多种循环类型:

- **For-each 循环**: 遍历数组或映射
- **While 循环**: 条件循环
- **Times 循环**: 固定次数循环

**示例**:

```yaml
steps:
  - id: process_items
    name: 处理每个项目
    loop:
      type: for_each
      items: ${vars.items}  # ["item1", "item2", "item3"]
      var: item
      index_var: idx
    prompt: |
      处理项目: ${loop.item}
      索引: ${loop.index}
```

### 6. 状态检测器

**文件位置**: `internal/status/detector.go`

实时监控代理状态:

- **空闲检测**: 基于活动阈值(默认 5 秒)
- **错误状态**: 通过模式匹配识别错误
- **输出预览**: 提供调试信息

## 执行流程

### 1. 初始化阶段

```
1. 加载工作流文件(YAML/TOML)
2. 验证工作流结构
3. 构建依赖图
4. 检测循环依赖
5. 初始化变量
6. 加载或创建执行状态
```

### 2. 执行阶段

```
主执行循环:
  while 还有未执行的步骤:
    1. 获取就绪步骤(GetReadySteps())
    2. 如果没有就绪步骤但还有未完成的步骤:
       - 报错: 死锁情况
    3. 对每个就绪步骤:
       a. 选择最佳代理
       b. 发送提示并监控执行
       c. 捕获输出
       d. 更新状态
       e. 标记步骤为已执行
```

### 3. 并行执行模式

```
并发控制:
- 使用信号量限制最大并发数(默认 8)
- 避免同一代理的并行使用
- 支持三种错误模式:
  * fail: 等待所有步骤完成,如果有失败则整体失败
  * fail_fast: 第一个失败时取消所有步骤
  * continue: 忽略错误继续执行
```

### 4. 状态持久化

```
状态存储:
- 位置: .ntm/pipelines/{workflow_id}_{timestamp}.json
- 原子写入: 使用临时文件和重命名
- 内容:
  * 步骤执行状态
  * 变量值
  * 输出结果
  * 错误信息
```

## 错误处理和重试

### 错误策略

在步骤和工作流级别配置:

```yaml
# 工作流级别
on_error: fail  # fail, continue, fail_fast

# 步骤级别
steps:
  - id: risky_step
    on_error: retry
    retry:
      max_attempts: 3
      backoff: exponential  # exponential, linear
      initial_delay: 1s
      max_delay: 30s
```

### 重试逻辑

1. **指数退避**: `delay = min(initial * 2^attempt, max_delay)`
2. **线性退避**: `delay = initial * attempt`
3. **最大延迟**: 防止无限增长的延迟

## 输出解析

支持多种输出解析方式:

```yaml
output:
  # 1. JSON 解析
  format: json
  path: ${steps.step_id.output}

  # 2. YAML 解析
  format: yaml

  # 3. 正则提取
  format: regex
  regex: 'Result: (\d+)'

  # 4. 键值对
  format: key_value
  separator: '='
```

## 配置文件示例

### 基本工作流

```yaml
id: example_workflow
name: 示例工作流
description: 展示基本的 pipeline 功能

vars:
  project_name: "my-project"
  items: ["a", "b", "c"]

steps:
  - id: setup
    name: 初始化
    agent: claude
    prompt: |
      初始化项目 {{.vars.project_name}}
    output:
      format: json
      path: ${steps.setup.output}

  - id: process
    name: 处理数据
    agent: codex
    depends_on: [setup]
    loop:
      type: for_each
      items: ${vars.items}
      var: item
    prompt: |
      处理项目: ${vars.project_name}
      当前项: ${loop.item}
      设置: ${steps.setup.output.config}

  - id: summarize
    name: 汇总结果
    agent: claude
    depends_on: [process]
    prompt: |
      请汇总以下处理结果:
      ${steps.process.output}
```

### 并行工作流

```yaml
id: parallel_example
name: 并行处理示例

steps:
  - id: fetch_data
    name: 获取数据
    agent: claude
    prompt: 获取需要处理的数据
    output:
      format: json

  - id: parallel_processing
    name: 并行处理
    depends_on: [fetch_data]
    parallel:
      steps:
        - id: process_a
          agent: codex
          prompt: 处理 A 类型: ${steps.fetch_data.output.type_a}

        - id: process_b
          agent: codex
          prompt: 处理 B 类型: ${steps.fetch_data.output.type_b}

        - id: process_c
          agent: gemini
          prompt: 处理 C 类型: ${steps.fetch_data.output.type_c}
      on_error: fail_fast  # 任一失败立即终止

  - id: combine
    name: 合并结果
    depends_on: [parallel_processing]
    prompt: |
      合并以下结果:
      A: ${steps.process_a.output}
      B: ${steps.process_b.output}
      C: ${steps.process_c.output}
```

## 性能优化

### 1. 并行执行

- 自动识别可并行步骤
- 可配置并发限制
- 智能代理分配

### 2. 状态管理

- 增量状态更新
- 压缩状态文件
- 定期清理旧状态

### 3. 缓存机制

- 代理响应缓存
- 变量解析缓存
- 依赖解析缓存

### 4. 连接池

- 复用代理连接
- 健康检查缓存
- 负载均衡优化

## 监控和调试

### 状态跟踪

```bash
# 查看执行状态
ntm pipeline status {workflow_id}

# 实时日志
ntm pipeline logs {workflow_id} --follow

# 性能指标
ntm pipeline metrics {workflow_id}
```

### 调试功能

- **详细日志**: 记录每个决策点
- **状态快照**: 定期保存完整状态
- **执行跟踪**: 步骤执行时间线
- **代理选择**: 记录选择理由

## 最佳实践

### 1. 工作流设计

- 保持步骤粒度适中
- 合理使用依赖关系
- 利用并行执行
- 添加适当的超时

### 2. 错误处理

- 为关键步骤添加重试
- 使用 fail_fast 快速失败
- 记录详细的错误信息
- 考虑降级策略

### 3. 性能优化

- 批量处理数据
- 减少不必要的变量传递
- 使用代理亲和性
- 监控资源使用

### 4. 可维护性

- 添加清晰的描述
- 使用有意义的步骤 ID
- 文档化工作流
- 版本控制工作流文件

## 扩展性

### 1. 自定义代理类型

```go
// 添加新的代理类型
func init() {
    RegisterAgentType("custom", &CustomAgent{})
}
```

### 2. 自定义输出解析器

```go
// 实现 OutputParser 接口
type CustomParser struct{}

func (p *CustomParser) Parse(data string) (interface{}, error) {
    // 自定义解析逻辑
}
```

### 3. 自定义路由策略

```go
// 实现 RoutingStrategy 接口
type CustomStrategy struct{}

func (s *CustomStrategy) SelectAgent(agents []*AgentState) *AgentState {
    // 自定义选择逻辑
}
```

## 总结

NTM Pipeline 是一个功能强大的工作流编排系统,具有以下特点:

1. **灵活性**: 支持复杂的依赖关系和执行模式
2. **可靠性**: 完善的错误处理和状态管理
3. **可扩展性**: 模块化的架构设计
4. **高性能**: 智能的并行执行和负载均衡
5. **易用性**: 声明式配置和丰富的调试工具

通过理解这些实现原理,用户可以更好地设计和优化他们的工作流,充分利用系统的强大功能。
