package main

import (
	"eino-learn/lesson1"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// 版本1，仅使用聊天模型进行对话
	lesson1.Agent()
}
