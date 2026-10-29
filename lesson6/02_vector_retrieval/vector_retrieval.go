// 语义空间向量检索
package vectorretrieval

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/schema"
)

func VectorRetrieval() {
	ctx := context.Background()
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		BaseURL: os.Getenv("EMBEDDING_URL"),
		Model:   os.Getenv("EMBEDDING_MODEL"),
		APIKey:  os.Getenv("EMBEDDING_API_KEY"),
	})
	if err != nil {
		log.Fatalf("new embedder err: %v", err)
	}

	// 获取预先准备的文档向量
	docs := []*schema.Document{
		{
			Content: "Go 语言由 Google 开发，2009年发布。它是一门静态类型、编译型语言。",
			MetaData: map[string]any{
				"source": "go_basics",
				"type":   "intro",
			},
		},
		{
			Content: "Go 的并发模型基于 CSP（通信顺序进程）。主要通过 goroutine 和 channel 实现。",
			MetaData: map[string]any{
				"source": "go_concurrency",
				"type":   "advanced",
			},
		},
		{
			Content: "Go 有着非常强大的并发能力，通过go关键词开启goroutine，每个goroutine的栈初始只有2kb，非常轻量。",
			MetaData: map[string]any{
				"source": "goroutine",
				"type":   "advanced",
			},
		},
		{
			Content: "Eino 是用 Go 开发的 AI 应用框架，提供了丰富的组件和灵活的编排能力。",
			MetaData: map[string]any{
				"source": "eino_intro",
				"type":   "framework",
			},
		},
	}

	// 将文档的内容抽取出来
	contents := make([]string, len(docs))
	for i, doc := range docs {
		contents[i] = doc.Content
	}

	// 转换为向量
	vectors, err := embedder.EmbedStrings(ctx, contents)
	if err != nil {
		log.Fatalf("embed strins err: %v", err)
	}

	// -------------进行语义检索-------------------

	query := "Go 语言的并发是如何实现的？"
	fmt.Printf("=== 检索: %s ===\n", query)

	// 先将用户提问转换为向量
	queryVec, err := embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		log.Fatalf("查询向量化失败: %v", err)
	}

	// 用于计算相似度并排序
	type docScore struct {
		doc   *schema.Document
		score float64
	}
	// 遍历文档对应的向量，挨个计算向量相似度
	scores := make([]docScore, 0)
	for i, doc := range docs {
		score := cosineSimilarity(queryVec[0], vectors[i])
		scores = append(scores, docScore{doc, score})
	}
	// 按相似度从高到低排序，取 Top 2
	sort.Slice(scores, func(i, j int) bool { return scores[i].score > scores[j].score })
	topK := 2

	// 获取检索到的文档
	retrievedDocs := make([]*schema.Document, topK)
	for i := 0; i < topK; i++ {
		retrievedDocs[i] = scores[i].doc
	}

	// 显示检索结果
	fmt.Printf("✅ 检索到 %d 个相关文档:\n\n", len(retrievedDocs))
	for i, doc := range retrievedDocs {
		fmt.Printf("文档 %d:\n", i+1)
		fmt.Printf("  内容: %s\n", doc.Content)
		fmt.Printf("  元数据: %v\n\n", doc.MetaData)
	}
}

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
