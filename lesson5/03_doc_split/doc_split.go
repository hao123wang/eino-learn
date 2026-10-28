package docsplit

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/schema"
)

// 文档拆分
func SplitDoc() {
	ctx := context.Background()

	// 创建长文档
	longText := `
	Go 语言是 Google 开发的一门编程语言。它具有以下特点：

	1. 并发支持：Go 内置了 goroutine 和 channel，让并发编程变得简单。
	2. 静态类型：在编译时进行类型检查，减少运行时错误。
	3. 快速编译：编译速度快，提升开发效率。
	4. 简洁语法：语法简单清晰，易于学习和维护。
	5. 标准库丰富：提供了完善的标准库，覆盖常见需求。

	Go 广泛应用于云计算、微服务、DevOps 等领域。许多知名公司如 Google、Uber、Docker 等都使用 Go 语言开发核心服务。

	Go 语言的并发模型基于 CSP（Communicating Sequential Processes）理论，通过 goroutine 和 channel 实现轻量级并发。
	每个 goroutine 只占用几 KB 的内存，可以轻松创建成千上万个 goroutine。

	Go 语言的工具链包括编译器、格式化工具、测试工具等，都集成在 go 命令中，使用非常方便。
	`

	// 切分
	chunks := recurSplit(ctx, longText)

	// 打印切分结果
	fmt.Printf("原文长度: %d 字符\n", len([]rune(longText)))
	fmt.Printf("分割成 %d 块\n\n", len(chunks))

	for i, chunk := range chunks {
		fmt.Printf("=== 块 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", chunk.Content)
		fmt.Printf("长度: %d 字符\n", len([]rune(chunk.Content)))
		fmt.Printf("元数据: %v\n\n", chunk.MetaData)
	}
}

// 手动切分文本内容
func manualSplit(s string) []*schema.Document {
	// 设置文本块大小和重叠范围
	chunkSize := 100
	chunkOverlap := 20

	textRunes := []rune(s)

	var chunks []*schema.Document
	for i := 0; i < len(textRunes); i += chunkSize - chunkOverlap {
		end := i + chunkSize
		end = min(end, len(textRunes))

		chunk := &schema.Document{
			Content: string(textRunes[i:end]),
			MetaData: map[string]any{
				"source": "go_intro.txt",
				"chunk":  len(chunks) + 1,
			},
		}
		chunks = append(chunks, chunk)
		if end >= len(textRunes) {
			break
		}
	}
	return chunks
}

// 递归切分文档
func recurSplit(ctx context.Context, s string) []*schema.Document {
	// 创建递归切分器
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   500,                           // 每个块最多500字符
		OverlapSize: 50,                            // 块之间重叠50个字符
		Separators:  []string{"\n", ".", "?", "!"}, // 分隔符列表，不设置的话默认就是这四种
		KeepType:    recursive.KeepTypeNone,        // 每个块丢弃分隔符
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建*schema.Documet文档
	doc := &schema.Document{
		Content: s,
		MetaData: map[string]any{
			"source": "go_intro.txt",
			"title":  "go语言简介",
		},
	}

	// 对文档进行切割
	chunks, err := splitter.Transform(ctx, []*schema.Document{doc})
	if err != nil {
		log.Fatal(err)
	}
	return chunks
}

// 按 #、##、### 多级标题切分markdown文件
func titleSplitMarkdown(ctx context.Context, s string) []*schema.Document {
	// 创建markdown切分器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		// 切割器在切分过程中，会根据markdown标题提取标题内容，填充到chunk的元数据中。
		// 元数据的键：Headers中进行配置，元数据的值：切分器从markdown文档中提取到的标题内容
		Headers: map[string]string{
			"#":   "title",      // 一级标题
			"##":  "section",    // 二级标题
			"###": "subsection", // 三级标题
		},
		TrimHeaders: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	doc := &schema.Document{
		Content: s,
		MetaData: map[string]any{
			"source": "go_intro.md",
		},
	}

	chunks, err := splitter.Transform(ctx, []*schema.Document{doc})
	if err != nil {
		log.Fatal(err)
	}
	return chunks
}
