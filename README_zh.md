# AG-UI Go 包

完整的 Agent User Interaction Protocol (AG-UI) Go 实现，用于构建具有流式对话、工具调用和状态管理的 AI 应用程序。

## 快速开始

### 安装

```bash
go get github.com/WuKongIM/WuKongIM/pkg/ag-ui
```

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    
    agui "github.com/WuKongIM/WuKongIM/pkg/ag-ui"
)

func main() {
    // 创建用户消息
    message := agui.NewUserMessage("msg_1", "你好，你能帮助我什么？", "user_123")
    
    // 编码为 JSON
    data, err := agui.EncodeMessage(message)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("编码后的消息: %s\n", string(data))
    
    // 解码回来
    decoded, err := agui.DecodeMessageFromBytes(data)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("来自 %s 的消息: %s\n", decoded.GetRole(), decoded.GetID())
}
```

## 核心概念

### 事件 (Events)

事件表示 AG-UI 中的实时交互。所有事件都实现 `Event` 接口：

```go
// 创建生命周期事件
runStarted := agui.NewRunStartedEvent("thread_1", "run_1")
runFinished := agui.NewRunFinishedEvent("thread_1", "run_1", nil)

// 创建流式文本消息事件
textStart := agui.NewTextMessageStartEvent("msg_1")
textContent := agui.NewTextMessageContentEvent("msg_1", "你好")
textEnd := agui.NewTextMessageEndEvent("msg_1")

// 编码事件
data, err := agui.EncodeEvent(textContent)
```

### 消息 (Messages)

消息表示对话历史。所有消息都实现 `Message` 接口：

```go
// 不同类型的消息
systemMsg := agui.NewSystemMessage("msg_1", "你是一个有用的助手。", "")
userMsg := agui.NewUserMessage("msg_2", "天气怎么样？", "user_123")
assistantMsg := agui.NewAssistantMessage("msg_3", "让我为你查看一下。", "assistant", nil)

// 带工具调用的助手消息
toolCall := agui.ToolCall{
    ID:   "tool_1",
    Type: agui.ToolCallTypeFunction,
    Function: agui.FunctionCall{
        Name:      "get_weather",
        Arguments: `{"location": "北京"}`,
    },
}
assistantWithTools := agui.NewAssistantMessage("msg_4", "", "assistant", []agui.ToolCall{toolCall})
```

### 工具调用 (Tool Calls)

使用流式事件处理工具执行：

```go
// 工具调用流程
toolStart := agui.NewToolCallStartEvent("tool_1", "get_weather", "msg_4")
toolArgs := agui.NewToolCallArgsEvent("tool_1", `{"location": "北京"}`)
toolEnd := agui.NewToolCallEndEvent("tool_1")
toolResult := agui.NewToolCallResultEvent("msg_5", "tool_1", "晴天，22°C")
```

## 流式事件

### 基本流式处理

```go
package main

import (
    "bytes"
    "fmt"
    "log"
    
    agui "github.com/WuKongIM/WuKongIM/pkg/ag-ui"
)

func streamingExample() {
    // 创建事件流
    events := []agui.Event{
        agui.NewRunStartedEvent("thread_1", "run_1"),
        agui.NewTextMessageStartEvent("msg_1"),
        agui.NewTextMessageContentEvent("msg_1", "你好"),
        agui.NewTextMessageContentEvent("msg_1", "世界！"),
        agui.NewTextMessageEndEvent("msg_1"),
        agui.NewRunFinishedEvent("thread_1", "run_1", nil),
    }
    
    // 编码为流
    var buf bytes.Buffer
    for _, event := range events {
        data, err := agui.EncodeEvent(event)
        if err != nil {
            log.Fatal(err)
        }
        buf.Write(data)
        buf.WriteString("\n")
    }
    
    // 解码流
    decoder := agui.NewStreamDecoder(&buf)
    eventChan, errorChan := decoder.DecodeEvents()
    
    for {
        select {
        case event, ok := <-eventChan:
            if !ok {
                return
            }
            fmt.Printf("收到事件: %s\n", event.GetType())
            
        case err := <-errorChan:
            if err != nil {
                log.Printf("流错误: %v", err)
                return
            }
        }
    }
}
```

### 实时消息构建

```go
func buildStreamingMessage() {
    var messageContent strings.Builder
    
    decoder := agui.NewStreamDecoder(reader)
    eventChan, errorChan := decoder.DecodeEvents()
    
    for {
        select {
        case event := <-eventChan:
            switch e := event.(type) {
            case *agui.TextMessageStartEvent:
                fmt.Printf("消息开始: %s\n", e.MessageID)
                messageContent.Reset()
                
            case *agui.TextMessageContentEvent:
                messageContent.WriteString(e.Delta)
                fmt.Printf("当前内容: %s\n", messageContent.String())
                
            case *agui.TextMessageEndEvent:
                fmt.Printf("最终消息: %s\n", messageContent.String())
                return
            }
            
        case err := <-errorChan:
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}
```

## Agent 输入和验证

### 创建 RunAgentInput

```go
func createAgentInput() *agui.RunAgentInput {
    // 定义工具
    searchTool := agui.Tool{
        Name:        "search",
        Description: "搜索信息",
        Parameters: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "query": map[string]interface{}{
                    "type":        "string",
                    "description": "搜索查询",
                },
            },
            "required": []string{"query"},
        },
    }
    
    // 创建消息
    messages := []agui.Message{
        agui.NewSystemMessage("msg_1", "你是一个有用的助手。", ""),
        agui.NewUserMessage("msg_2", "搜索 Go 教程", "user_123"),
    }
    
    // 创建上下文
    context := []agui.Context{
        {
            Description: "用户偏好",
            Value:       "偏好简洁的解释",
        },
    }
    
    // 创建完整输入
    input := &agui.RunAgentInput{
        ThreadID:       agui.GenerateThreadID(),
        RunID:          agui.GenerateRunID(),
        State:          map[string]interface{}{"step": 1},
        Messages:       messages,
        Tools:          []agui.Tool{searchTool},
        Context:        context,
        ForwardedProps: map[string]interface{}{"version": "1.0"},
    }
    
    // 验证输入
    if err := input.Validate(); err != nil {
        log.Fatalf("无效输入: %v", err)
    }
    
    return input
}
```

### 验证示例

```go
func validationExamples() {
    // 有效事件
    event := &agui.RunStartedEvent{
        BaseEvent: agui.BaseEvent{Type: agui.EventTypeRunStarted},
        ThreadID:  "thread_123",
        RunID:     "run_456",
    }
    
    if err := event.Validate(); err != nil {
        fmt.Printf("验证失败: %v\n", err)
    } else {
        fmt.Println("事件有效")
    }
    
    // 无效事件（缺少必需字段）
    invalidEvent := &agui.TextMessageContentEvent{
        BaseEvent: agui.BaseEvent{Type: agui.EventTypeTextMessageContent},
        MessageID: "msg_123",
        Delta:     "", // 空 delta 无效
    }
    
    if err := invalidEvent.Validate(); err != nil {
        fmt.Printf("预期验证错误: %v\n", err)
    }
}
```

## 状态管理

### 状态快照和增量

```go
func stateManagement() {
    // 初始状态
    initialState := map[string]interface{}{
        "user_id":     "user_123",
        "step_count":  1,
        "preferences": map[string]string{"theme": "dark"},
    }
    
    // 创建状态快照
    snapshot := agui.NewStateSnapshotEvent(initialState)
    
    // 创建状态增量（JSON Patch 操作）
    delta := []interface{}{
        map[string]interface{}{
            "op":    "replace",
            "path":  "/step_count",
            "value": 2,
        },
        map[string]interface{}{
            "op":    "replace",
            "path":  "/preferences/theme",
            "value": "light",
        },
    }
    
    deltaEvent := agui.NewStateDeltaEvent(delta)
    
    // 编码并发送
    snapshotData, _ := agui.EncodeEvent(snapshot)
    deltaData, _ := agui.EncodeEvent(deltaEvent)
    
    fmt.Printf("状态快照: %s\n", string(snapshotData))
    fmt.Printf("状态增量: %s\n", string(deltaData))
}
```

## API 参考

### 核心类型

- **EventType**: 所有事件类型的枚举（17 种类型）
- **Role**: 消息角色的枚举（developer, system, assistant, user, tool）
- **ToolCallType**: 工具调用类型的枚举（function）

### 主要接口

```go
type Event interface {
    GetType() EventType
    GetTimestamp() *int64
    GetRawEvent() interface{}
    Validate() error
    EventTypeName() string
}

type Message interface {
    GetID() string
    GetRole() Role
    GetName() string
    Validate() error
    MessageType() string
}
```

### 编码函数

```go
// 事件编码
func EncodeEvent(event Event) ([]byte, error)
func DecodeEventFromBytes(data []byte) (Event, error)

// 消息编码
func EncodeMessage(message Message) ([]byte, error)
func DecodeMessageFromBytes(data []byte) (Message, error)

// 流处理
func NewStreamDecoder(r io.Reader) *StreamDecoder
func (s *StreamDecoder) DecodeEvents() (<-chan Event, <-chan error)
func (s *StreamDecoder) DecodeMessages() (<-chan Message, <-chan error)
```

### 工厂函数

```go
// 事件工厂
func NewRunStartedEvent(threadID, runID string) *RunStartedEvent
func NewTextMessageContentEvent(messageID, delta string) *TextMessageContentEvent
func NewToolCallStartEvent(toolCallID, toolCallName, parentMessageID string) *ToolCallStartEvent

// 消息工厂
func NewUserMessage(id, content, name string) *UserMessage
func NewAssistantMessage(id, content, name string, toolCalls []ToolCall) *AssistantMessage

// ID 生成器
func GenerateMessageID() string
func GenerateRunID() string
func GenerateThreadID() string
func GenerateToolCallID() string
```

## 常见模式

### 完整对话流程

```go
func conversationFlow() {
    // 1. 开始 agent 运行
    runStarted := agui.NewRunStartedEvent("thread_1", "run_1")

    // 2. 流式助手响应
    msgStart := agui.NewTextMessageStartEvent("msg_1")
    msgContent1 := agui.NewTextMessageContentEvent("msg_1", "我会帮助你解决这个问题。")
    msgContent2 := agui.NewTextMessageContentEvent("msg_1", "让我搜索一些信息。")
    msgEnd := agui.NewTextMessageEndEvent("msg_1")

    // 3. 进行工具调用
    toolStart := agui.NewToolCallStartEvent("tool_1", "search", "msg_1")
    toolArgs := agui.NewToolCallArgsEvent("tool_1", `{"query": "Go 教程"}`)
    toolEnd := agui.NewToolCallEndEvent("tool_1")

    // 4. 工具结果
    toolResult := agui.NewToolCallResultEvent("msg_2", "tool_1", "找到 10 个 Go 教程")

    // 5. 最终响应
    finalStart := agui.NewTextMessageStartEvent("msg_3")
    finalContent := agui.NewTextMessageContentEvent("msg_3", "这里有一些很棒的 Go 教程...")
    finalEnd := agui.NewTextMessageEndEvent("msg_3")

    // 6. 完成运行
    runFinished := agui.NewRunFinishedEvent("thread_1", "run_1", map[string]int{"tutorials_found": 10})

    events := []agui.Event{
        runStarted, msgStart, msgContent1, msgContent2, msgEnd,
        toolStart, toolArgs, toolEnd, toolResult,
        finalStart, finalContent, finalEnd, runFinished,
    }

    // 处理事件
    for _, event := range events {
        data, _ := agui.EncodeEvent(event)
        fmt.Printf("事件: %s\n", string(data))
    }
}
```

### 错误处理

```go
func errorHandling() {
    // 带验证的编码
    event := agui.NewTextMessageContentEvent("msg_1", "你好")
    data, err := agui.EncodeEvent(event)
    if err != nil {
        switch {
        case errors.Is(err, agui.ErrValidationFailed):
            log.Printf("验证错误: %v", err)
        case errors.Is(err, agui.ErrMarshalFailed):
            log.Printf("编码错误: %v", err)
        default:
            log.Printf("未知错误: %v", err)
        }
        return
    }

    // 带错误处理的解码
    decoded, err := agui.DecodeEventFromBytes(data)
    if err != nil {
        switch {
        case errors.Is(err, agui.ErrUnmarshalFailed):
            log.Printf("解码错误: %v", err)
        case errors.Is(err, agui.ErrInvalidEventType):
            log.Printf("无效事件类型: %v", err)
        default:
            log.Printf("未知错误: %v", err)
        }
        return
    }

    fmt.Printf("成功处理事件: %s\n", decoded.GetType())
}
```

## 最佳实践

### 性能提示

- **重用解码器**: 每个连接创建一个 `StreamDecoder`
- **批量操作**: 尽可能批量处理多个事件
- **提前验证**: 编码前使用验证方法
- **处理错误**: 始终检查编码/解码错误

### 内存管理

```go
// 好的做法：重用解码器
decoder := agui.NewStreamDecoder(conn)
eventChan, errorChan := decoder.DecodeEvents()

// 好的做法：批量处理
var events []agui.Event
for i := 0; i < batchSize; i++ {
    select {
    case event := <-eventChan:
        events = append(events, event)
    case <-time.After(timeout):
        break
    }
}
processBatch(events)
```

### 线程安全

```go
// 安全：每个 goroutine 有自己的解码器
func handleConnection(conn net.Conn) {
    decoder := agui.NewStreamDecoder(conn)
    eventChan, errorChan := decoder.DecodeEvents()

    for {
        select {
        case event := <-eventChan:
            // 安全地处理事件
        case err := <-errorChan:
            // 处理错误
            return
        }
    }
}
```

## 故障排除

### 常见问题

**1. 验证错误**
```go
// 问题：必需字段为空
event := &agui.TextMessageContentEvent{
    BaseEvent: agui.BaseEvent{Type: agui.EventTypeTextMessageContent},
    MessageID: "msg_1",
    Delta:     "", // 这会导致验证失败
}

// 解决方案：确保所有必需字段都已设置
event.Delta = "你好世界"
```

**2. 类型不匹配**
```go
// 问题：错误的事件类型
event := &agui.RunStartedEvent{
    BaseEvent: agui.BaseEvent{Type: agui.EventTypeRunFinished}, // 错误类型
    ThreadID:  "thread_1",
    RunID:     "run_1",
}

// 解决方案：使用正确的类型
event.BaseEvent.Type = agui.EventTypeRunStarted
```

**3. JSON 解析错误**
```go
// 问题：函数参数中的无效 JSON
toolCall := agui.ToolCall{
    Function: agui.FunctionCall{
        Name:      "search",
        Arguments: `{invalid json}`, // 这会导致验证失败
    },
}

// 解决方案：使用有效的 JSON
toolCall.Function.Arguments = `{"query": "搜索词"}`
```

### 调试技巧

- 编码前使用 `Validate()` 方法检查数据
- 使用 `errors.Is()` 检查错误类型以进行特定处理
- 为流处理启用详细日志记录
- 先用小数据集进行测试

## 测试

运行测试套件：

```bash
cd pkg/ag-ui
go test -v
```

运行特定测试：

```bash
go test -v -run TestEventEncoding
go test -v -run TestStreamDecoding
```

## 贡献

1. 遵循 Go 约定和最佳实践
2. 为新功能添加测试
3. 更新 API 变更的文档
4. 根据 AG-UI 协议规范进行验证

## 许可证

此实现遵循 AG-UI 协议规范，与其他 AG-UI 实现兼容。

有关 AG-UI 协议的更多信息：https://docs.ag-ui.com/
```
```
