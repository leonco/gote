package gote

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type IncludeNode struct {
	name  string
	flags []Flag
}

func parseIncludeNode(spec string) IncludeNode {
	split := strings.Split(spec, ":")
	return IncludeNode{split[0], parseModifiers(split[1:])}
}

func (in IncludeNode) Evaluate(dict *TemplateDictionary, context *TemplateLoaderContext, w io.Writer) error {
	filename, ok := dict.Get(in.name)
	if !ok {
		panic("The template identifier for included section " + in.name + " is not set!")
	}

	inclTmpl, err := context.Loader.GetTemplate(filename)
	if err != nil {
		panic(err)
	}
	inclTmpl.context = context
	var preW io.Writer
	var buf = &bytes.Buffer{}
	if len(in.flags) > 0 {
		preW = w
		w = buf
	}
	dicts := dict.ChildDicts(in.name)

	if len(dicts) == 0 {
		inclTmpl.Render(dict, w)
	} else {
		for _, d := range dicts {
			inclTmpl.Render(d, w)
		}
	}

	if preW != nil {
		w = preW
		results := buf.Bytes()
		_, err := w.Write(applyModifiers(results, in.flags))
		if err != nil {
			return err
		}
	}
	return nil
}

func (i IncludeNode) String() string {
	return fmt.Sprintf("INCLUDE(%s, flag=%d)", i.name, i.flags)
}
