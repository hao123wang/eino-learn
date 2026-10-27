// graph 编排 允许有分支和环
package graphmodel

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// 用于保存graph的局部状态
type state struct {
	Messages []*schema.Message
}

type TimeParams struct {
	Format string `json:"format"`
}

// 获取时间函数出参
type TimeResult struct {
	CurrentTime string `json:"current_time"`
}

// 获取时间函数
func GetCurrentTime(ctx context.Context, params *TimeParams) (*TimeResult, error) {
	now := time.Now()
	var result string
	switch params.Format {
	case "date":
		result = now.Format("2006-01-02")
	case "time":
		result = now.Format("15:04:05")
	default:
		result = now.Format("2006-01-02 15:04:05")
	}
	return &TimeResult{CurrentTime: result}, nil
}

func ChatModel() {
	ctx := context.Background()

	// 1. prompt 模板
	chatTempl := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是{role}"),
		schema.UserMessage("{question}"),
	)

	// 2. ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 3. tool
	bizTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"format": {Type: schema.String, Desc: "获取当前时间", Required: true},
		}),
	}, GetCurrentTime)
	// 工具信息
	toolInfo, err := bizTool.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// 工具节点
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{bizTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 4. toolModel
	toolModel, err := chatModel.WithTools([]*schema.ToolInfo{toolInfo})
	if err != nil {
		log.Fatal(err)
	}

	// 5. graph
	const (
		chatTemplKey = "chatTemplNode"
		toolModelKey = "toolModelNode"
		toolsNodeKey = "toolsNode"
	)
	graph := compose.NewGraph[map[string]any, *schema.Message](
		compose.WithGenLocalState(func(ctx context.Context) *state {
			return &state{}
		}),
	)

	// 添加节点
	must(graph.AddChatTemplateNode(chatTemplKey, chatTempl))
	must(graph.AddChatModelNode(toolModelKey, toolModel,
		compose.WithStatePreHandler(func(ctx context.Context, in []*schema.Message, s *state) ([]*schema.Message, error) {
			s.Messages = append(s.Messages, in...)
			return s.Messages, nil
		}),
	))

	must(graph.AddToolsNode(toolsNodeKey, toolsNode,
		compose.WithStatePreHandler(func(ctx context.Context, in *schema.Message, s *state) (*schema.Message, error) {
			// 工具节点要的是带 tool_calls 的那条 assistant 消息，历史只记在 state 里
			s.Messages = append(s.Messages, in)
			return in, nil
		}),
	))
	// 添加边
	must(graph.AddEdge(compose.START, chatTemplKey)) // 起始节点为prompt
	must(graph.AddEdge(chatTemplKey, toolModelKey))  // prompt -> toolChatModel
	// 工具结果回到同一个模型。pre handler 会把 tool 消息追加进 state，再把完整历史交给模型
	must(graph.AddEdge(toolsNodeKey, toolModelKey))

	// 模型节点后添加分支，根据模型输出决定是否调用工具
	// graphBranch第一个参数是函数，根据输入的内容进行判断是否调用工具，决定走向下一个节点。 第二个参数列出所有可能走的节点
	must(graph.AddBranch(toolModelKey, compose.NewGraphBranch(
		func(ctx context.Context, in *schema.Message) (string, error) {
			if len(in.ToolCalls) > 0 {
				fmt.Println("  [branch] 模型请求调用工具 → tools")
				return toolsNodeKey, nil
			}
			fmt.Println("  [branch] 模型直接回答 → END")
			return compose.END, nil
		},
		map[string]bool{toolsNodeKey: true, compose.END: true},
	)))

	// 准备输入，并进行编译
	// 工具节点指回模型，图里有环，必须限制最大步数
	input := map[string]any{
		"role":     ` 一名ai编程助手。只做与编程有关的事情，不进行其他方面的知识回答。每次用户提问，都要在响应内容前面加上date类型的当前时间`,
		"question": "请你介绍一下golang语言的特点, 限制在300字以内",
	}
	runnable, err := graph.Compile(ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		log.Fatal(err)
	}
	// 运行
	ouput, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ouput.Content)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
