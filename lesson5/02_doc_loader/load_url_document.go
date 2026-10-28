package loaddocument

import (
	"context"
	"fmt"
	"log"

	urlloader "github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
)

// 加载url文档
func LoadURLDocument() {
	ctx := context.Background()

	// 创建文档加载器
	loader, err := urlloader.NewLoader(ctx, &urlloader.LoaderConfig{})
	if err != nil {
		log.Fatal(err)
	}

	// 加载文档
	docs, err := loader.Load(ctx, document.Source{
		URI: "https://example.com", // Web URL
	})
	if err != nil {
		log.Fatal(err)
	}
	// 打印加载的文档
	fmt.Printf("加载了 %d 个文档\n", len(docs))
	for i, doc := range docs {
		fmt.Printf("\n=== 文档 %d ===\n", i+1)
		// 只打印前 200 个字符，避免输出过长
		content := doc.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		fmt.Printf("内容预览: %s\n", content)
		fmt.Printf("内容长度: %d 字符\n", len(doc.Content))
		fmt.Printf("元数据: %v\n", doc.MetaData)
	}
}
