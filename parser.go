package gote

import (
	"bytes"
	"errors"
	"log"
	"regexp"
)

var (
	OPEN_SQUIGGLE    = regexp.QuoteMeta("{{")
	CLOSE_SQUIGGLE   = regexp.QuoteMeta("}}")
	VARIABLE_RE      = "([a-zA-Z_]+(:[a-zA-Z]+)*)*"
	RE_OPEN_SECTION  = regexp.MustCompile("(?im:" + OPEN_SQUIGGLE + "#([a-zA-Z_]+)" + CLOSE_SQUIGGLE + ")")
	RE_CLOSE_SECTION = regexp.MustCompile("(?im:" + OPEN_SQUIGGLE + "/([a-zA-Z_]+)" + CLOSE_SQUIGGLE + ")")
	RE_VARIABLE      = regexp.MustCompile("(?im:" + OPEN_SQUIGGLE + VARIABLE_RE + CLOSE_SQUIGGLE + ")")
	RE_INCLUDE       = regexp.MustCompile("(?im:" + OPEN_SQUIGGLE + ">" + VARIABLE_RE + CLOSE_SQUIGGLE + ")")
)

type NodeType int

const (
	OPEN_SECTION NodeType = iota
	CLOSE_SECTION
	VARIABLE
	TEXT_NODE
	INCLUDE_SECTION
)

type TemplateParser interface {
	Parse(template []byte) ([]TemplateNode, error)
}

type CTemplateParser struct {
}

func NewCTemplateParser() TemplateParser {
	return &CTemplateParser{}
}

func (t *CTemplateParser) Parse(input []byte) (ns []TemplateNode, err error) {
	for len(input) > 0 {
		log.Printf("looking ahead at '%s'\n", input)
		switch next(input) {
		case TEXT_NODE:
			input, ns = handleTextNode(input, ns)
		case OPEN_SECTION:
			input, ns = handleOpenSection(input, ns)
		case CLOSE_SECTION:
			input, ns = handleCloseSection(input, ns)
		case INCLUDE_SECTION:
			input, ns = handleInclude(input, ns)
		case VARIABLE:
			input, ns = handleVariable(input, ns)
		default:
			err = errors.New("Internal error parsing template.")
		}
	}
	return
}

func handleTextNode(input []byte, nodes []TemplateNode) (r []byte, ns []TemplateNode) {
	nextBraces := bytes.Index(input, []byte("{{"))
	var text []byte
	if nextBraces < 0 {
		text = input
	} else {
		text = input[:nextBraces]
		r = input[nextBraces:]
	}
	if len(text) > 0 {
		ns = append(nodes, TextNode{text})
	}
	return
}

func handleInclude(input []byte, nodes []TemplateNode) (r []byte, ns []TemplateNode) {
	c, r := consume(RE_INCLUDE, input)
	ns = append(nodes, parseIncludeNode(string(c)))
	return
}

func handleOpenSection(input []byte, nodes []TemplateNode) (r []byte, ns []TemplateNode) {
	c, r := consume(RE_OPEN_SECTION, input)
	ns = append(nodes, OpenSection(string(c)))
	return
}

func handleCloseSection(input []byte, nodes []TemplateNode) (r []byte, ns []TemplateNode) {
	c, r := consume(RE_CLOSE_SECTION, input)
	ns = append(nodes, CloseSection(string(c)))
	return
}

func handleVariable(input []byte, nodes []TemplateNode) (r []byte, ns []TemplateNode) {
	c, r := consume(RE_VARIABLE, input)
	ns = append(nodes, parseVNode(string(c)))
	return
}

func next(input []byte) NodeType {
	if bytes.HasPrefix(input, []byte("{{#")) {
		return OPEN_SECTION
	} else if bytes.HasPrefix(input, []byte("{{/")) {
		return CLOSE_SECTION
	} else if bytes.HasPrefix(input, []byte("{{>")) {
		return INCLUDE_SECTION
	} else if bytes.HasPrefix(input, []byte("{{")) {
		return VARIABLE
	} else {
		return TEXT_NODE
	}
}

func consume(p *regexp.Regexp, input []byte) (c, r []byte) {
	matches := p.FindSubmatch(input)

	l := len(matches)
	if l < 2 {
		return nil, input
	}
	return matches[1], input[len(matches[0]):]
}
