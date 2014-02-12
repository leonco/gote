package gote

import (
	"fmt"
	"io"
)

//
// Represents a literal string.
//
type TextNode struct {
	text []byte
}

func (tnode TextNode) Evaluate(dict *TemplateDictionary, context *TemplateLoaderContext, w io.Writer) error {
	_, err := w.Write(tnode.text)
	return err
}

func (t TextNode) String() string {
	return fmt.Sprintf("TEXT(%s)", t.text)
}
