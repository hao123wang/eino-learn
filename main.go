package main

import (
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// 版本1，仅使用聊天模型进行对话
	// lesson1.Agent()

	// 版本2，使用chain/graph/workflow编排，实现agent
	// chain.Agent()
	// graph.Agent()
	// workflow.Agent()
}
