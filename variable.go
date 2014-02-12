package gote

import (
	"fmt"
	"io"
	"strings"
)

//
// Represents a node whose output is defined by a value from the TemplateDictionary.
//
// This supports both {{PLAIN}} variables as well as one with {{MODIFERS:j}}.
//
type VariableNode struct {
	variable  string
	flags []Flag
}

func parseVNode(spec string) VariableNode {
	ss := strings.Split(spec, ":")
	return VariableNode{ss[0], parseModifiers(ss[1:])}
}

func (vnode VariableNode) Evaluate(dict *TemplateDictionary, context *TemplateLoaderContext, w io.Writer) error {
	t, ok := dict.Get(vnode.variable)
	if ok {
		w.Write([]byte(t))
	}
	return nil
}

func (v VariableNode) String() string {
	return fmt.Sprintf("VARIABLE(%s, flag=%d)", v.variable, v.flags)
}
