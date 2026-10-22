package graph

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

	// ------------------------创建输入节点---------------------------------

	// 创建聊天模型
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
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
		schema.MessagesPlaceholder("history", true),
		schema.UserMessage("{user_query}"),
	)

	// 创建工具
	stockInfoTool := utils.NewTool(&schema.ToolInfo{
		Name: "play_info",
		Desc: "根据用户提供的股票代码或是企业名称，查询股票的相关信息, 并根据信息去做一个详细的投资分析",
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

	// 创建工具节点
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{stockInfoTool},
	})
	if err != nil {
		log.Fatalf("failed to create tools node: %v", err)
	}

	// 创建一个lambda节点
	lambdaNode := compose.InvokableLambda(biztools.ToolToTextLambda)

	// 给每个节点一个key
	chatTemplKey := "chat_template"
	toolsNodeKey := "tools_node"
	lambdaToolToTextKey := "lambda_tool_to_text"
	chatModelKey := "chat_model"

	// 分支：模型发起了工具调用就去工具节点，否则说明已得到最终回答，结束
	toolCallBranch := compose.NewGraphBranch(
		func(ctx context.Context, msg *schema.Message) (string, error) {
			if len(msg.ToolCalls) > 0 {
				return toolsNodeKey, nil
			}
			return compose.END, nil
		},
		map[string]bool{toolsNodeKey: true, compose.END: true},
	)

	// ------------------------graph编排---------------------------------
	// 各个节点之间，需要用边进行连接，边是有向的，从起点到终点
	// graph可以使用分支，如果模型发起了工具调用就去工具节点，否则说明已得到最终回答，结束
	graph := compose.NewGraph[map[string]any, *schema.Message]()

	// 添加节点
	graph.AddChatTemplateNode(chatTemplKey, chatTempl)
	graph.AddChatModelNode(chatModelKey, chatModel)
	graph.AddToolsNode(toolsNodeKey, toolsNode)
	graph.AddLambdaNode(lambdaToolToTextKey, lambdaNode)

	// 添加边
	graph.AddEdge(compose.START, chatTemplKey)
	graph.AddEdge(chatTemplKey, chatModelKey)
	graph.AddEdge(toolsNodeKey, lambdaToolToTextKey)
	graph.AddEdge(lambdaToolToTextKey, chatModelKey)

	// 添加分支
	graph.AddBranch(chatModelKey, toolCallBranch)

	// ------------------------编译运行---------------------------------
	// 限制最大运行步数，避免模型反复调用工具导致死循环
	runnable, err := graph.Compile(ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		log.Fatalf("compile graph failed: %v", err)
	}

	output, err := runnable.Invoke(ctx, map[string]any{
		"user_query": "我想查询一下腾讯的股票信息",
	})
	if err != nil {
		log.Fatalf("invoke graph failed: %v", err)
	}
	println("=====================思考内容====================")
	if output.ReasoningContent != "" {
		println(output.ReasoningContent)
	}
	println("=========================================")
	if output.Content != "" {
		println(output.Content)
	}
}
