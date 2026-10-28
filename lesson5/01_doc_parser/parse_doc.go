// 文档解析器，从io.reader中将文本内容解析为[]*schema.Document类型
package docparser

import (
	"context"
	"fmt"
	"log"
	"os"

	htmlparser "github.com/cloudwego/eino-ext/components/document/parser/html"
	pdfparser "github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document/parser"
)

func ParseDoc() {
	ctx := context.Background()

	// 创建text文本解析器
	textParser := parser.TextParser{}
	// 创建html文件解析器
	htmlParser, _ := htmlparser.NewParser(ctx, &htmlparser.Config{
		Selector: &htmlparser.BodySelector, // 只提取body部分的内容
	})
	// 创建pdf文件解析器
	pdfParser, _ := pdfparser.NewPDFParser(ctx, &pdfparser.Config{
		ToPages: false, // false：合并所有页面为一个文档   true：每页一个文档
	})

	textFile, err := os.Open("lesson5/testdata/sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer textFile.Close()

	htmlFile, err := os.Open("lesson5/testdata/sample.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlFile.Close()

	pdfFile, err := os.Open("lesson5/testdata/sample.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer pdfFile.Close()

	textDocs, err := textParser.Parse(ctx, textFile)
	if err != nil {
		log.Fatal(err)
	}
	htmlDocs, err := htmlParser.Parse(ctx, htmlFile)
	if err != nil {
		log.Fatal(err)
	}
	pdfDocs, err := pdfParser.Parse(ctx, pdfFile)
	if err != nil {
		log.Fatal(err)
	}

	// 打印解析结果
	fmt.Printf("解析了 %d 个 text 文档\n", len(textDocs))
	for i, doc := range textDocs {
		fmt.Printf("\n=== 文档 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", doc.Content)
		fmt.Printf("元数据: %v\n", doc.MetaData)
	}

	fmt.Printf("解析了 %d 个 html 文档\n", len(htmlDocs))
	for i, doc := range htmlDocs {
		fmt.Printf("\n=== 文档 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", doc.Content)
		fmt.Printf("元数据: %v\n", doc.MetaData)
	}

	fmt.Printf("解析了 %d 个 pdf 文档\n", len(pdfDocs))
	for i, doc := range pdfDocs {
		fmt.Printf("\n=== 文档 %d ===\n", i+1)
		fmt.Printf("内容: %s\n", doc.Content)
		fmt.Printf("元数据: %v\n", doc.MetaData)
	}
}
