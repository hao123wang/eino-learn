package main

import (
	workflow "eino-learn/lesson4/03_workflow"
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
	// tools.ChatModel()

	// lesson4
	// chainmodel.ChatModel()
	// graphmodel.ChatModel()
	workflow.ChatModel()
}
