package main

import (
	"eino-learn/lesson3/tools"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("load environment variable err: %v", err)
	}
	// lesson1
	// chat.ChatModel()
	// multiturn.ChatModel()
	// param.ChatModel()

	// lesson2
	// chattemplate.ChatModel()
	// chathistory.ChatModel()

	// lesson3
	tools.ChatModel()
}
