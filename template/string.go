package template

import (
	"strings"

	"golang.org/x/net/html"
)

type stringTemplate struct {
	text string
}

func (t *stringTemplate) GetFragment(context *html.Node) []*html.Node {
	nodes, _ := html.ParseFragment(strings.NewReader(t.text), context)
	return nodes
}

func (t *stringTemplate) MarshalText() (text []byte, err error) {
	return []byte(t.text), nil
}

func String(text string) Template {
	return &stringTemplate{text}
}
