package lib

import (
	"io"

	"github.com/ebadfd/jira_sucks/pkg/jirawiki"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func renderHeading(w io.Writer, p *ast.Heading, entering bool) {
	if entering {
		io.WriteString(w, "<b>")
	} else {
		io.WriteString(w, "</b> <br>")
	}
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if heading, ok := node.(*ast.Heading); ok {
		renderHeading(w, heading, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: myRenderHook,
	}
	return html.NewRenderer(opts)
}

func MarkdownToHtml(md string) []byte {
	res := jirawiki.Parse(string(md))

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(
		[]byte(res),
	)

	renderer := newCustomizedRender()
	return markdown.Render(doc, renderer)
}
