// chain 编排示例
package chain

import (
	"context"
	"eino-learn/biztools"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func Agent() {
	ctx := context.Background()

	// 创建聊天模型
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   os.Getenv("MODEL_ID"),
		BaseURL: os.Getenv("BASE_URL"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}

	// ChatTemplate 节点
	sysPrompt := `你是一个资深的金融投资分析师，你将根据用户提供的股票代码或是企业名称，查询股票的相关信息,
		并根据信息去做一个详细的投资分析，从行业背景、企业的财报，一些关键的参考指标，来给用户提供一个详细的投资分析。`
	chatTempl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(sysPrompt),
		// 作为占位符，在prompt中预留历史消息插入位置, bool表示是否允许接受空的消息列表（用户可能首次对话）
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("{user_query}"), // 用户输入的问题
	)

	// 创建工具
	stockInfoTool := utils.NewTool(&schema.ToolInfo{
		Name: "play_info",
		Desc: "根据用户提供的股票代码或是企业名称，查询股票的相关信息, 并根据信息去做一个详细的投资分析",
		// 工具入参, 与工具函数入参对应
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"stock_code": {
				Type: schema.String,
				Desc: "股票代码",
			},
			"company_name": {
				Type: schema.String,
				Desc: "企业名称",
			},
		}),
	}, biztools.GetStockInfo)

	// 工具绑定给模型
	toolInfo, err := stockInfoTool.Info(ctx)
	if err != nil {
		log.Fatalf("failed to get tool info: %v", err)
	}
	if err := chatModel.BindTools([]*schema.ToolInfo{toolInfo}); err != nil {
		log.Fatalf("failed to bind tools: %v", err)
	}

	// 工具节点
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{stockInfoTool},
	})
	if err != nil {
		log.Fatalf("failed to create tools node: %v", err)
	}

	// 创建一个lambda节点
	lambdaToolToText := compose.TransformableLambda(biztools.ToolToTextLambda)

	// chain编排 NewChain[输入类型, 输出类型]
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.AppendChatTemplate(chatTempl)  // 聊天模板：将用户输入的问题，转换为模板化的聊天消息
	chain.AppendChatModel(chatModel)     // 模型节点：根据输入的聊天消息，生成响应消息
	chain.AppendToolsNode(toolsNode)     // 工具节点：负责执行工具调用，返回工具执行结果
	chain.AppendLambda(lambdaToolToText) // lambda节点：将工具调用结果转换为文本
	chain.AppendChatModel(chatModel)     // 模型节点：根据工具调用结果，生成响应消息

	// 执行chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("compile chain failed: %v", err)
	}

	// 执行chain
	output, err := runnable.Invoke(ctx, map[string]any{
		"user_query": "我想查询一下腾讯控股的股票信息",
	})
	if err != nil {
		log.Fatalf("invoke chain failed: %v", err)
	}

	println("=====================思考内容====================")
	if output.ReasoningContent != "" {
		println(output.ReasoningContent)
	}
	println("=========================================")
	if output.Content != "" {
		println(output.Content)
	}

	// 可选：token 用量（如果 deepseek 的 response meta 有 usage）
	if output.ResponseMeta != nil && output.ResponseMeta.Usage != nil {
		println("\nToken 使用统计:")
		println("  输入 Token:", int(output.ResponseMeta.Usage.PromptTokens))
		println("  输出 Token:", int(output.ResponseMeta.Usage.CompletionTokens))
		println("  总计 Token:", int(output.ResponseMeta.Usage.TotalTokens))
	}
}
