package main

import (
	"eino-learn/lesson2/graph"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// 版本1，仅使用聊天模型进行对话
	// lesson1.Agent()

	// 版本2，使用chain编排，实现agent
	graph.Agent()
}
