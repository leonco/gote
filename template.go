package gote

import (
	"fmt"
	"io"
	"log"
)

var (
	parser TemplateParser = NewCTemplateParser()
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
	return parseTemplate(parser, content)
}

func parseTemplate(parser TemplateParser, content []byte) (*Template, error) {
	nodes, err := parser.Parse(content)

	if err != nil {
		return nil, err
	}

	return newTemplate(nodes), nil
}

func (t *Template) Render(td *TemplateDictionary, w io.Writer) {
	render(t.context, t.tmpl, td, w)
}

func render(context *TemplateLoaderContext, nodes []TemplateNode, dict *TemplateDictionary, w io.Writer) {
	for _, node := range nodes {
		switch nnode := node.(type) {
		case SectionNode:
			log.Println(nnode)
		default:
			node.Evaluate(dict, context, w)
		}
	}
}

func handleSection(context *TemplateLoaderContext, nodes []TemplateNode, dict *TemplateDictionary, openTagIdx int, w io.Writer) (int, error) {
	sn := nodes[openTagIdx].(SectionNode)
	snName := sn.name
	p := openTagIdx + 1

	otherSections := 0

	for ; p < len(nodes); p++ {
		node := nodes[p]

		if tag, ok := node.(SectionNode); ok {
			if tag.isOpenSectionTag() {
				otherSections += 1
			}

			if tag.isCloseSectionTag() {
				if tag.name == snName {
					break
				} else if otherSections == 0 {
					return p, fmt.Errorf("mismatched close tag: expecting a close tag for %s, "+
						"but got close tag for %s", snName, tag.name)
				} else {
					otherSections--
				}
			}
		}
	}

	if p == len(nodes) {
		return p, fmt.Errorf("missing close tag for %s ", snName)
	}

	if dict.isHidenSection(snName) {
		log.Printf("Skipping section %s because it is hidden", snName)
		return p, nil
	}

	subdicts := dict.ChildDicts(snName)

	if len(subdicts) == 0 {
		render(context, nodes[openTagIdx+1:p], dict, w)
	} else {
		for _, td := range subdicts {
			render(context, nodes[openTagIdx+1:p], td, w)
		}
	}

	return p, nil
}
