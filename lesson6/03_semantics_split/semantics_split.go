// 按语义切分文档，在前一段和后一段语义相差大的地方进行切分
package semanticssplit

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/schema"
)

func SemanticsSplit() {
	ctx := context.Background()
	embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		BaseURL: os.Getenv("EMBEDDING_URL"),
		Model:   os.Getenv("EMBEDDING_MODEL"),
		APIKey:  os.Getenv("EMBEDDING_API_KEY"),
	})
	if err != nil {
		log.Fatalf("new embedder err: %v", err)
	}

	// 文档切分流程
	// 根据Separators，按顺序各切一轮，分隔符留在左边那段，获取多个chunk。如hello.world!你好啊  切分后是hell.  world!
	// BufferSize决定当前chunk前后拼接几个chunk。例如第i段送去embedding的文本是 前n段 + 当前段 + 后n段
	splitter, err := semantic.NewSplitter(ctx, &semantic.Config{
		Embedding:    embedder,                      // 用于计算语义相似度的 Embedding 模型
		Percentile:   0.7,                           // 分割阈值（0-1 之间，默认 0.9）
		BufferSize:   1,                             // 每个块前后拼接的块数量（可选）
		MinChunkSize: 50,                            // 最小块大小（可选）
		Separators:   []string{"\n", ".", "?", "!"}, // 分隔符列表（可选，默认值）
	})
	if err != nil {
		log.Fatalf("创建分割器失败: %v", err)
	}

	// 创建文档
	doc := &schema.Document{
		Content: `
		第一段：Go 语言是 Google 开发的一门编程语言。它具有并发支持、静态类型、快速编译等特点。

		第二段：Go 语言广泛应用于云计算、微服务、DevOps 等领域。许多知名公司都使用 Go 语言开发核心服务。

		第三段：Go 语言的并发模型基于 CSP 理论，通过 goroutine 和 channel 实现轻量级并发。

		第四段：Go 语言的工具链包括编译器、格式化工具、测试工具等，都集成在 go 命令中。
		`,
		MetaData: map[string]any{
			"source": "go_intro.txt",
		},
	}
	
	// 分割文档
	chunks, err := splitter.Transform(ctx, []*schema.Document{doc})
	if err != nil {
		log.Fatalf("分割文档失败: %v", err)
	}

	// 打印分割结果
	fmt.Printf("原文长度: %d 字符\n", len([]rune(doc.Content)))
	fmt.Printf("分割成 %d 块\n\n", len(chunks))

	for i, chunk := range chunks {
		fmt.Printf("=== 块 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", chunk.Content)
		fmt.Printf("长度: %d 字符\n\n", len([]rune(chunk.Content)))
	}
}
