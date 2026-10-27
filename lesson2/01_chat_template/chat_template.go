package chattemplate

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

	// 创建ChatTemplate
	chatTempl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		schema.UserMessage("{question}"),
	)
	// 准备变量
	variables := map[string]any{
		"role":     "资深程序员",
		"question": "请你用golang写出令牌桶算法代码",
	}
	// 格式化消息
	inMessage, err := chatTempl.Format(ctx, variables)
	if err != nil {
		log.Fatalf("ChatTemplate format err: %v\n", err)
	}

	// 创建模型
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	response, err := chatModel.Generate(ctx, inMessage)
	if err != nil {
		log.Fatalf("生成失败: %v", err)
	}

	fmt.Printf("\nAI 回答:\n%s\n", response.Content)
}
