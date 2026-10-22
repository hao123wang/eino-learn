package biztools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

// 工具是参与编排流程的节点，将业务逻辑封装为工具，由model来决策是否调用工具

// 工具入参结构体
// json tag 必须与 ToolInfo.ParamsOneOf 中的参数名一致，否则模型传来的参数无法反序列化，字段会是空值
type StockInfoInput struct {
	StockCode   string `json:"stock_code"`   // 股票代码
	CompanyName string `json:"company_name"` // 企业名称
}

// 工具出参结构体
type StockInfoOutput struct {
	Info string `json:"info"` // 股票信息
}

// 工具执行函数
func GetStockInfo(ctx context.Context, input *StockInfoInput) (*StockInfoOutput, error) {
	// 联网查询股票信息, api key调用的model，只是一个推理模型，不具备联网搜索、rag等功能
	return &StockInfoOutput{
		Info: fmt.Sprintf("股票代码: %s, 企业名称: %s, 股票信息: 暂无实时数据（模拟）", input.StockCode, input.CompanyName),
	}, nil
}

// Lambda是给编排流程中调用的节点，将业务逻辑封装为lambda，参与编排流程的执行

// 将工具调用结果转换为文本。
// ChatModel 的入参是 []*schema.Message，所以这里的输出类型也必须是消息列表，不能是单条 *schema.Message。
func ToolToTextLambda(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
	// 从 toolsNode 输出里找 role=tool 的消息，把它的 content 摘出来
	var toolContents []string
	for _, m := range msgs {
		if m != nil && m.Role == "tool" {
			toolContents = append(toolContents, m.Content)
		}
	}
	// 没找到 tool 内容也没关系，给一个兜底文本，避免空
	text := "工具未返回有效信息。"
	if len(toolContents) > 0 {
		// 多个工具结果就拼起来
		text = "工具返回的用户信息如下：\n" + "- " + toolContents[0]
		for i := 1; i < len(toolContents); i++ {
			text += "\n- " + toolContents[i]
		}
	}

	return []*schema.Message{schema.UserMessage(text)}, nil
}
