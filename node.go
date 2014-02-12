package gote

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"net/url"
	"text/template"
	"strings"
)



type Flag int

const (
	H Flag = iota // HTML
	X             // XML
	J             // JavaScript string literal
	U             // URL escaped (%xx)
	B             // turns \n into <br/>
)

var flagTable = map[string]Flag{"H": H, "X": X, "J": J, "U": U, "B": B}

type Range struct {
	start, end, skipTo int
}

type TemplateNode interface {
	Evaluate(dict *TemplateDictionary, context *TemplateLoaderContext, w io.Writer) error
}

func parseModifiers(split []string) []Flag {
	flags := make([]Flag, 10)
	i := 0
	for _, fs := range split {
		f, ok := flagTable[strings.ToUpper(fs)]
		if ok {
			flags[i] = f
			i += 1
		}
	}
	return flags[:i]
}

func applyModifiers(input []byte, modifiers []Flag) []byte {
	for _, m := range modifiers {
		switch m {
		case H:
			input = htmlEscape(input)
		case J:
			input = jsEscape(input)
		case X:
			input = xmlEscape(input)
		case U:
			input = urlEncode(input)
		case B:
			input = newlinesToBreaks(input)
		}
	}
	return input
}

func jsEscape(input []byte) []byte {
	var b bytes.Buffer
	template.JSEscape(&b, input)
	return b.Bytes()
}

func htmlEscape(input []byte) []byte {
	var b bytes.Buffer
	template.HTMLEscape(&b, input)
	return b.Bytes()
}

func xmlEscape(input []byte) []byte {
	var b bytes.Buffer
	xml.EscapeText(&b, input)
	return b.Bytes()
}

func urlEncode(input []byte) []byte {
	s := url.QueryEscape(string(input))
	return []byte(s)
}

func newlinesToBreaks(input []byte) []byte {
	return bytes.Replace(input, []byte("\n"), []byte("<br />"), -1)
}

func TestNode() {
	log.Print("testing node.go")
}
