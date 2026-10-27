package param

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func ChatModel() {
	ctx := context.Background()

	// 创建聊天模型时，设置相关参数
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),

		// 模型参数设置
		Temperature: 1.0,
		TopP:        0.95,
		MaxTokens:   800,
	})
	if err != nil {
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	inMessage := []*schema.Message{
		schema.SystemMessage("你是一个知识渊博的百度百科ai助手，名字叫小助, 你每次的回答都简短到300字以内"),
		schema.UserMessage("你好，请介绍一下河南地区比较有名的美食"),
	}

	// 调用模型生成响应
	// 这里只对本次调用生效：温度降到 0.2、最多 60 个 token，Config 里的默认值不受影响
	response, err := chatModel.Generate(ctx, inMessage,
		model.WithTemperature(1.0),
		model.WithMaxTokens(1000),
	)
	if err != nil {
		log.Fatalf("生成响应失败: %v", err)
	}
	fmt.Println(response.Content)
}
