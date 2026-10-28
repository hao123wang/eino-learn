// 文档加载器, 将文档加载为[]*schema.Document类型
package loaddocument

import (
	"context"
	"fmt"
	"log"

	fileloader "github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
)

// 加载本地文档
func LoadLocal() {
	ctx := context.Background()

	// 创建解析器
	textParser := parser.TextParser{}
	// htmlParser, _ := htmlparser.NewParser(ctx, &htmlparser.Config{
	// 	Selector: &htmlparser.BodySelector, // 只提取body部分的内容
	// })
	// pdfParser, _ := pdfparser.NewPDFParser(ctx, &pdfparser.Config{
	// 	ToPages: false, // false：合并所有页面为一个文档   true：每页一个文档
	// })

	// 创建一个本地文件加载器
	loader, err := fileloader.NewFileLoader(ctx, &fileloader.FileLoaderConfig{
		UseNameAsID: true,       // 是否使用文件名作为文档id
		Parser:      textParser, // 可以指定不同的解析器，否则默认会解析所有内容
	})
	if err != nil {
		log.Fatal(err)
	}

	// 从本地文件中加载文档
	docs, err := loader.Load(ctx, document.Source{
		URI: "lesson5/testdata/sample.html",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("加载了 %d 个文档\n", len(docs))
	for i, doc := range docs {
		fmt.Printf("\n=== 文档 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", doc.Content)
		fmt.Printf("元数据: %v\n", doc.MetaData)
	}
}
