package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// 获取时间函数入参
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

	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	// 通过utils.NewTool声明一个工具
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"format": {
				Type:     schema.String,
				Desc:     "时间格式: date(日期), time(时间), datetime(日期时间)",
				Required: false,
			},
		}),
	}, GetCurrentTime)

	// 通过实现接口声明的工具
	calculatorTool := &CalculatorTool{}

	// 收集工具和工具信息
	tools := []tool.InvokableTool{timeTool, calculatorTool}
	var toolInfos []*schema.ToolInfo
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			log.Fatalf("获取工具信息失败: %v", err)
		}
		toolInfos = append(toolInfos, info)
	}

	// 创建ToolsNode，负责执行工具
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{timeTool, calculatorTool},
	})
	if err != nil {
		log.Fatalf("创建 ToolsNode 失败: %v", err)
	}

	testCases := []string{
		"现在几点了？",
		"帮我计算 10.5 + 20 等于多少",
		"你好，请介绍一下你自己",
	}
	for i, question := range testCases {
		fmt.Printf("\n========== 测试 %d ==========\n", i+1)
		fmt.Printf("问题: %s\n", question)

		messages := []*schema.Message{schema.UserMessage(question)}

		// 第 1 步：model.WithTools 只对这一次调用生效
		response, err := chatModel.Generate(ctx, messages, model.WithTools(toolInfos))
		if err != nil {
			log.Printf("生成失败: %v", err)
			continue
		}

		// 第 2 步：没有 ToolCalls 说明模型选择直接回答
		if len(response.ToolCalls) == 0 {
			fmt.Printf("  ✓ 直接回答: %s\n", response.Content)
			continue
		}

		// 第 3 步：调用工具执行
		// ToolsNode 的输入是带 ToolCalls 的 assistant 消息，输出是每个 ToolCall 对应的一条 role=tool 消息	}
		toolResults, err := toolsNode.Invoke(ctx, response)
		if err != nil {
			log.Printf("  工具执行失败: %v", err)
			continue
		}

		// 将本轮工具调用，添加到上下文
		messages = append(messages, response)
		messages = append(messages, toolResults...)
		final, err := chatModel.Generate(ctx, messages)
		if err != nil {
			log.Printf("  最终回答生成失败: %v", err)
			continue
		}
		fmt.Printf("  ✓ 最终回答: %s\n", final.Content)
	}
}
