package gote

import (
	"fmt"
	"io"
)

//
// Implementation of a {{#SECTION_NODE}} and the paired {{/SECTION_NODE}}.
//
type SectionNode struct {
	name string
	t    NodeType
}

func (sn SectionNode) Evaluate(dict *TemplateDictionary, context *TemplateLoaderContext, w io.Writer) error {
	return nil
}

func (sn SectionNode) isOpenSectionTag() bool {
	return sn.t == OPEN_SECTION
}

func (sn SectionNode) isCloseSectionTag() bool {
	return sn.t == CLOSE_SECTION
}

func OpenSection(name string) SectionNode {
	return SectionNode{name, OPEN_SECTION}
}

func CloseSection(name string) SectionNode {
	return SectionNode{name, CLOSE_SECTION}
}

func (s SectionNode) String() string {
	return fmt.Sprintf("SECTION(%s, %d)", s.name, s.t)
}
