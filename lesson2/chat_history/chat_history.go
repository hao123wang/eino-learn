package chathistory

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func ChatModel() {
	ctx := context.Background()

	// 创建模型
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	// 创建ChatTemplate
	chatTempl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}, 回答控制在两句话以内"),
		schema.MessagesPlaceholder("history", false),
		schema.UserMessage("{question}"),
	)

	questions := []string{
		"我叫小明，正在学习go语言。",
		"推荐我下一步学什么?",
		"我叫什么名字?",
	}
	var history []*schema.Message
	for _, q := range questions {
		params := map[string]any{
			"role":     "编程助手",
			"history":  history,
			"question": q,
		}
		inMessage, err := chatTempl.Format(ctx, params)
		if err != nil {
			log.Fatalf("格式化失败: %v", err)
		}
		response, err := chatModel.Generate(ctx, inMessage)
		if err != nil {
			log.Fatalf("生成失败: %v", err)
		}
		fmt.Printf("\nAI 回答:\n%s\n", response.Content)
		history = append(history, schema.UserMessage(q), response)
	}
}
