// chain 编排
package chainmodel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

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
		log.Fatalf("创建 ChatModel 实例失败: %v", err)
	}

	// 3. chain
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.
		AppendChatTemplate(chatTempl).
		AppendChatModel(chatModel)

	// 编译
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("chain compile err: %v", err)
	}

	// 准备输入内容
	input := map[string]any{
		"role":     "一名资深的心理咨询师，正在进行角色扮演来帮助一名抑郁的小女孩",
		"question": "请你向小女孩介绍一下自己",
	}
	// 运行
	// output, err := runnable.Invoke(ctx, input)
	output, err := runnable.Stream(ctx, input)
	if err != nil {
		log.Fatalf("runnable.Stream err: %v", err)
	}

	// 读取流式输出
	for {
		chunk, err := output.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("output.Recv err: %v", err)
		}
		fmt.Print(chunk.Content)
	}
	fmt.Println()
}
