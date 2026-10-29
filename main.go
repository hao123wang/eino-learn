package main

import (
	semanticssplit "eino-learn/lesson6/03_semantics_split"
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
	// workflow.ChatModel()

	// lesson5
	// docparser.ParseDoc()
	// docparser.ExtParse()
	// docloader.LoadURLDocument()
	// docsplit.SplitDoc()

	// lesson6
	// textembedding.TextEmbedding()
	// vectorretrieval.VectorRetrieval()
	semanticssplit.SemanticsSplit()
}
