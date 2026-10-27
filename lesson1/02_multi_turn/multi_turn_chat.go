// 多轮对话 + 流式响应
package multiturn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
	"github.com/peterh/liner"
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

	// 用来保存对话消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个知识渊博的百度百科ai助手，名字叫小助, 你每次的回答都简短到300字以内"),
	}

	line := liner.NewLiner()
	defer line.Close()

	for {
		// 读取终端输入
		input, err := line.Prompt("input>")
		if err != nil {
			log.Printf("line read err: %v", err)
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		if input == "" {
			continue
		}

		// 将输入添加到messages
		messages = append(messages, schema.UserMessage(input))

		// 获取流式响应
		stream, err := chatModel.Stream(ctx, messages)
		if err != nil {
			log.Fatalf("chat model get stream err: %v", err)
		}
		defer stream.Close()

		var chunks []*schema.Message
		for {
			// 接收响应
			chunk, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("接收失败: %v", err)
			}
			fmt.Print(chunk.Content)
			chunks = append(chunks, chunk)
		}
		fmt.Println()

		// 合并流式输出
		outMessage, err := schema.ConcatMessages(chunks)
		if err != nil {
			log.Printf("合并消息失败: %v", err)
			continue
		}
		// 将完整的 Assistant 消息加入历史
		messages = append(messages, outMessage)
	}
}
