package gote

import (
	"io"
	"log"
)

var (
	Parser TemplateParser = NewCTemplateParser()
)

type Template struct {
	tmpl    []TemplateNode
	loader  *TemplateLoader
	context *TemplateLoaderContext
}

func newTemplate(nodes []TemplateNode) *Template {
	return &Template{tmpl: nodes}
}

func Parse(content []byte) (*Template, error) {
	return nil, nil
}

func ParseTemplate(parser TemplateParser, content []byte) (*Template, error) {
	return nil, nil
}

func (t *Template) Render(td *TemplateDictionary, w io.Writer) {
	render(t.context, t.tmpl, td, w)
}

func render(context *TemplateLoaderContext, nodes []TemplateNode, dict *TemplateDictionary, w io.Writer) {
	for _, node := range nodes {
		switch nnode := node.(type) {
		case *SectionNode:
			log.Println(nnode)
		default:
			node.Evaluate(dict, context, w)
		}
	}

}
