package chat

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
)

func ChatModel() {
	ctx := context.Background()

	// 创建 ChatModel 实例
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	// 准备消息
	inMessage := []*schema.Message{
		schema.SystemMessage("你是一个知识渊博的百度百科ai助手，名字叫小助, 你每次的回答都简短到300字以内"),
		schema.UserMessage("你好，请介绍一下河南地区比较有名的美食"),
	}

	// 调用模型生成响应
	response, err := chatModel.Generate(ctx, inMessage)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}

	// 打印结果
	fmt.Printf("AI 响应: %s\n", response.Content)

	// 输出 token 使用情况：ResponseMeta 可能为 nil，取字段前必须判空
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Printf("\nToken 使用统计:\n")
		fmt.Printf("  输入 Token: %d\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf("  输出 Token: %d\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf("  总计 Token: %d\n", response.ResponseMeta.Usage.TotalTokens)
	}
}
