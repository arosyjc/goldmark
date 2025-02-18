package parsetest

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// return (ast.GoToNext, true) to tell html renderer to skip rendering this node
// (because you've rendered it)
func renderHookCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {

	// skip all nodes that are not CodeBlock nodes
	if _, ok := node.(*ast.CodeBlock); !ok {
		return ast.GoToNext, false
	}

	// custom rendering logic for ast.CodeBlock. By doing nothing it won't be
	// 写入自定义数据
	w.Write([]byte(`<span>Hello World</span>`))

	return ast.GoToNext, false
}

func TestParse2(t *testing.T) {

	inputFile := "markdown.md"           // 输入 md 文件
	inputTemplateFile := "template.html" // 模板文件
	outputFile := "markdown.html"        // 输出 html 文件

	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: renderHookCodeBlock, // 设置自定义处理钩子
	}
	renderer := html.NewRenderer(opts)

	// 打开 md 文件
	fs, err := os.Open(inputFile)
	if err != nil {
		return
	}

	// 读取 md 内容
	var md string
	scanner := bufio.NewScanner(fs)
	for scanner.Scan() {
		md += scanner.Text() + "\n"
	}
	if scanner.Err() != nil {
		log.Panicln(err)
	}

	// md to html
	html := markdown.ToHTML([]byte(md), nil, renderer)

	tp, err := template.ParseFiles(inputTemplateFile)
	if err != nil {
		return
	}

	// 写入输出文件
	f, err := os.Create(outputFile)
	if err != nil {
		return
	}
	tp.Execute(f, template.HTML(html))

	print(string(html))
}
