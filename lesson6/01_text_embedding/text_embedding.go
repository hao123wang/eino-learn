// 将文本向量化
package textembedding

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

func TextEmbedding() {
	ctx := context.Background()

	// 创建Embedding模型
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		BaseURL: os.Getenv("EMBEDDING_URL"),
		Model:   os.Getenv("EMBEDDING_MODEL"),
		APIKey:  os.Getenv("EMBEDDING_API_KEY"),
	})
	if err != nil {
		log.Fatalf("new bedder err: %v", err)
	}

	// 待向量化的文本
	texts := []string{
		"我喜欢吃香蕉",
		"我喜欢吃西瓜",
		"我喜欢吃菠萝",
		"我喜欢吃很多水果",
	}

	// 文本生成向量
	vectors, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		log.Fatalf("embedding err: %v", err)
	}

	// 打印结果
	for i, text := range texts {
		fmt.Printf("文本 %d: %s\n", i+1, text)
		fmt.Printf("  向量维度: %d\n", len(vectors[i]))
		fmt.Printf("  前5维: %v\n\n", vectors[i][:5])
	}

	// 计算相似度
	similarity12 := cosineSimilarity(vectors[0], vectors[1])
	similarity13 := cosineSimilarity(vectors[0], vectors[2])

	fmt.Printf("文本1 和 文本2 的相似度: %.4f\n", similarity12)
	fmt.Printf("文本1 和 文本3 的相似度: %.4f\n", similarity13) //
}

// 计算向量的余弦相似度
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dotProduct / (math.Sqrt(normA) + math.Sqrt(normB))
}
