package goldmark_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/testutil"
	"github.com/yuin/goldmark/text"
)

func TestAstParseJson(t *testing.T) {
	// body:=``
	body, _ := os.ReadFile("mark.md")
	md := New()
	n := md.Parser().Parse(text.NewReader(body))
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch node := n.(type) {
			case *ast.Blockquote:
				fmt.Println("BlockQuote", string(node.Lines().Value(body)))
			case *ast.FencedCodeBlock:
				fmt.Println("FencedCodeBlock", string(node.Lines().Value(body)))
			case *ast.HTMLBlock:
				fmt.Println("HTMLBlock", string(node.Lines().Value(body)))
			case *ast.Text:
				// fmt.Println("Text", string(node.Value(body)))
			case *ast.RawHTML:
				fmt.Println("RawHTML", string(node.Segments.Value(body)))
			case *ast.Paragraph:
				fmt.Println("Paragraph", string(node.Lines().Value(body)))
			default:
				// fmt.Println(node.Kind().String())
			}
		}
		return ast.WalkContinue, nil
	})
}
func TestASTBlockNodeText(t *testing.T) {
	var cases = []struct {
		Name   string
		Source string
		T1     string
		T2     string
		C      bool
	}{
		{
			Name: "AtxHeading",
			Source: `# l1

a

# l2`,
			T1: `l1`,
			T2: `l2`,
		},
		{
			Name: "SetextHeading",
			Source: `l1
l2
===============

a

l3
l4
==============`,
			T1: `l1
l2`,
			T2: `l3
l4`,
		},
		{
			Name: "CodeBlock",
			Source: `    l1
    l2

a

    l3
	l4`,
			T1: `l1
l2
`,
			T2: `l3
l4
`,
		},
		{
			Name: "FencedCodeBlock",
			Source: "```" + `
l1
l2
` + "```" + `

a

` + "```" + `
l3
l4`,
			T1: `l1
l2
`,
			T2: `l3
l4
`,
		},
		{
			Name: "Blockquote",
			Source: `> l1
> l2

a

> l3
> l4`,
			T1: `l1
l2`,
			T2: `l3
l4`,
		},
		{
			Name: "List",
			Source: `- l1
  l2

a

- l3
  l4`,
			T1: `l1
l2`,
			T2: `l3
l4`,
			C: true,
		},
		{
			Name: "HTMLBlock",
			Source: `<div>
l1
l2
</div>

a

<div>
l3
l4`,
			T1: `<div>
l1
l2
</div>
`,
			T2: `<div>
l3
l4`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := New()
			n := md.Parser().Parse(text.NewReader(s))
			// fmt.Println(n.ChildCount())
			// attributes := n.Attributes()
			ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					switch node := n.(type) {
					case *ast.Document:
						fmt.Println("Document", string(node.Lines().Value(s)))
					case *ast.Blockquote:
						fmt.Println("BlockQuote", string(node.Lines().Value(s)))
					case *ast.Text:
						fmt.Println("Text:", string(node.Value(s)))
					case *ast.FencedCodeBlock:
						fmt.Println("FencedCodeBlock", string(node.Lines().Value(s)))
					case *ast.Paragraph:
						fmt.Println("Paragraph", string(node.Lines().Value(s)))
					case *ast.List:
						fmt.Println("List", string(node.Lines().Value(s)))
					case *ast.ListItem:
						fmt.Println("ListItem", string(node.Lines().Value(s)))
					case *ast.TextBlock:
						fmt.Println("TextBlock", string(node.Lines().Value(s)))
					case *ast.CodeBlock:
					case *ast.String:
						fmt.Println(string(node.Value))
						// fmt.Println(v)
					case *ast.Heading:
						// fmt.Printf("Heading %d: %v\n", v.Level, v)
					default:
						// fmt.Println(v.Kind().String())
					}
				}
				return ast.WalkContinue, nil
			})
			// for n.HasChildren() {
			// 	// fmt.Println("属性", n)
			// 	switch v := n.(type) {
			// 	case *ast.Document:
			// 		fmt.Println("Document", v.Kind().String())
			// 	case *ast.Blockquote:
			// 		fmt.Println("BlockQuote")
			// 	case *ast.FencedCodeBlock:
			// 		fmt.Println("123")
			// 	case *ast.Paragraph:
			// 		fmt.Println("Paragraph", v.Lines().Value())
			// 	case *ast.List:
			// 		fmt.Println("List")
			// 	case *ast.ListItem:
			// 		fmt.Println("ListItem")
			// 	case *ast.TextBlock:
			// 		fmt.Println("TextBlock")
			// 	case *ast.CodeBlock:
			// 	case *ast.String:
			// 		// fmt.Println(string(v.Value))
			// 		// fmt.Println(v)
			// 	case *ast.Heading:
			// 		// fmt.Printf("Heading %d: %v\n", v.Level, v)
			// 	default:
			// 		// fmt.Println(v.Kind().String())
			// 	}
			// 	n = n.LastChild()
			// }
			// for idx, a := range attributes {
			// 	fmt.Println("属性", idx, a)
			// }
			// c1 := n.FirstChild()

			// c2 := c1.NextSibling().NextSibling()
			// if cs.C {
			// 	c1 = c1.FirstChild()
			// 	c2 = c2.FirstChild()
			// }
			// if !bytes.Equal(c1.Text(s), []byte(cs.T1)) { // nolint: staticcheck

			// 	t.Errorf("%s unmatch: %s", cs.Name, testutil.DiffPretty(c1.Text(s), []byte(cs.T1))) // nolint: staticcheck

			// }
			// if !bytes.Equal(c2.Text(s), []byte(cs.T2)) { // nolint: staticcheck

			// 	t.Errorf("%s(EOF) unmatch: %s", cs.Name, testutil.DiffPretty(c2.Text(s), []byte(cs.T2))) // nolint: staticcheck

			// }
		})
	}

}

func TestASTInlineNodeText(t *testing.T) {
	var cases = []struct {
		Name   string
		Source string
		T1     string
	}{
		{
			Name:   "CodeSpan",
			Source: "`c1`",
			T1:     `c1`,
		},
		{
			Name:   "Emphasis",
			Source: `*c1 **c2***`,
			T1:     `c1 c2`,
		},
		{
			Name:   "Link",
			Source: `[label](url)`,
			T1:     `label`,
		},
		{
			Name:   "AutoLink",
			Source: `<http://url>`,
			T1:     `http://url`,
		},
		{
			Name:   "RawHTML",
			Source: `<span>c1</span>`,
			T1:     `<span>`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := New()
			n := md.Parser().Parse(text.NewReader(s))
			c1 := n.FirstChild().FirstChild()
			if !bytes.Equal(c1.Text(s), []byte(cs.T1)) { // nolint: staticcheck
				t.Errorf("%s unmatch:\n%s", cs.Name, testutil.DiffPretty(c1.Text(s), []byte(cs.T1))) // nolint: staticcheck
			}
		})
	}

}
