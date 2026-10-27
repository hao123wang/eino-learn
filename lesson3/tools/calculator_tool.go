package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 声明一个结构体，实现tool.InvokableTool接口

// 工具函数的入参
type CalculatorParams struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// 工具函数的出参
type CalculatorResult struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

// CalculatorTool 计算器工具
type CalculatorTool struct{}

// 编译期检查：确认 CalculatorTool 实现了 tool.InvokableTool，方法签名写错会直接编译失败
var _ tool.InvokableTool = (*CalculatorTool)(nil)

// 实现工具接口的Info方法
func (t *CalculatorTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "calculator",
		Desc: "执行基本的数学计算（加、减、乘、除）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"operation": {
				Type:     "string",
				Desc:     "运算类型: add(加), subtract(减), multiply(乘), divide(除)",
				Required: true,
			},
			"a": {
				Type:     "number",
				Desc:     "第一个数字",
				Required: true,
			},
			"b": {
				Type:     "number",
				Desc:     "第二个数字",
				Required: true,
			},
		}),
	}, nil
}

// 实现工具接口的InvokableRun方法
func (t *CalculatorTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 1. 解析参数
	var params CalculatorParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
		return "", fmt.Errorf("解析参数失败: %w", err)
	}

	// 2. 执行计算
	var result float64
	switch params.Operation {
	case "add":
		result = params.A + params.B
	case "subtract":
		result = params.A - params.B
	case "multiply":
		result = params.A * params.B
	case "divide":
		// 业务错误放进返回的 JSON 而不是返回 error：模型能读到"除数不能为0"并换个说法回答用户；
		// 返回 error 则会让整个调用链失败
		if params.B == 0 {
			resultJSON, _ := json.Marshal(CalculatorResult{
				Error: "除数不能为0",
			})
			return string(resultJSON), nil
		}
		result = params.A / params.B
	default:
		resultJSON, _ := json.Marshal(CalculatorResult{
			Error: "不支持的运算类型",
		})
		return string(resultJSON), nil
	}

	// 3. 返回结果
	resultJSON, err := json.Marshal(CalculatorResult{
		Result: result,
	})
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}
