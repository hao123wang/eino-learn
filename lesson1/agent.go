// 版本1, 输入 -> Model -> 输出
package lesson1

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

func Agent() {
	// 构建上下文
	ctx := context.Background()

	// 创建ChatModel
	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		panic(err)
	}

	// 创建输出信息
	message := []*schema.Message{
		schema.SystemMessage("你是一个中文ai助手，名字叫小离, 始终用中文回答问题"),
		schema.UserMessage("你好，请介绍下自己"),
	}

	// 调用模型获取输出
	response, err := model.Generate(ctx, message)
	// 流式输出用 model.Stream()
	if err != nil {
		panic(err)
	}

	// 打印输出
	fmt.Println(response.Content)
	// 打印此次对话的token使用情况
	if response.ResponseMeta != nil && response.ResponseMeta.Usage != nil {
		fmt.Println("此次对话的token使用情况如下:")
		fmt.Printf("prompt_tokens: %d\n", response.ResponseMeta.Usage.PromptTokens)
		fmt.Printf("completion_tokens: %d\n", response.ResponseMeta.Usage.CompletionTokens)
		fmt.Printf("total_tokens: %d\n", response.ResponseMeta.Usage.TotalTokens)
	}
}
