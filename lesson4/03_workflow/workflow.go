// workflow 编排 有向无环，可以将上一节点的输出内容，映射到下一节点的输入字段
// chain/graph 中，下一个节点只能拿到上一个节点的完整输出
// workflow 中，一个节点可以从多个上游获取一部分字段，拼成自己的输入
//
//	AddInput("节点A")                                  整个输出直接作为输入
//	AddInput("节点A", compose.MapFields("x", "y"))     取 A 输出的 x 字段，放到本节点输入的 y 字段
//	AddInput("节点B", compose.ToField("z"))            把 B 的整个输出，放到本节点输入的 z 字段
package workflow

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type history struct {
	Prompt    []*schema.Message // system/user的消息
	Assistant *schema.Message   // 带tool_calls的assistant
	Tools     []*schema.Message // role=tool的消息
}

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

	// 1. prompt 模板
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage("{role}"),
		schema.UserMessage("{question}"),
	)

	// 2. ChatModel
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: os.Getenv("BASE_URL"),
		Model:   os.Getenv("MODEL_ID"),
		APIKey:  os.Getenv("API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 3. tool
	timeTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_current_time",
		Desc: "获取当前时间",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"format": {Type: schema.String, Desc: "指定时间输出格式", Required: true},
		}),
	}, GetCurrentTime)
	timeToolInfo, err := timeTool.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{timeTool},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 4. 带有tool的模型，由于workflow中不能使用循环，因此用带工具的模型和不带工具的模型进行编排
	toolModel, err := chatModel.WithTools([]*schema.ToolInfo{timeToolInfo})
	if err != nil {
		log.Fatal(err)
	}

	// 5. workflow
	wf := compose.NewWorkflow[map[string]any, *schema.Message]()
	wf.AddChatTemplateNode("prompt", chatTpl).AddInput(compose.START) // 起点
	wf.AddChatModelNode("toolModel", toolModel).AddInput("prompt")    // 接收完整的prompt节点输出
	wf.AddToolsNode("toolNode", toolsNode).AddInput("toolModel")      // 接收工具调用模型的输出（可能携带工具调用）

	// 包装一个lambda将历史对话添加到当前节点的输入
	wf.AddLambdaNode("merge", compose.InvokableLambda(
		func(ctx context.Context, in *history) ([]*schema.Message, error) {
			msgs := make([]*schema.Message, 0, len(in.Prompt)+1+len(in.Tools))
			msgs = append(msgs, in.Prompt...)
			msgs = append(msgs, in.Assistant)
			msgs = append(msgs, in.Tools...)
			return msgs, nil
		},
	)).
		AddInput("prompt", compose.ToField("Prompt")).
		AddInput("toolModel", compose.ToField("Assistant")).
		AddInput("toolNode", compose.ToField("Tools"))

	wf.AddChatModelNode("chatModel", chatModel).AddInput("merge")
	wf.End().AddInput("chatModel")

	// 编译
	runnable, err := wf.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// 运行
	output, err := runnable.Invoke(ctx, map[string]any{
		"role":     "一个ai编程助手, 只回答编程相关的知识内容，每次回答不得超过300字，且在回答时，在内容前加上当前时间",
		"question": "想成为一个云原生架构师, 详细的学习路线是什么",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output.Content)
}
