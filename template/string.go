package template

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type StringTemplate struct {
	text string
}

func (t *StringTemplate) Fragment(context *html.Node) []*html.Node {
	nodes, _ := html.ParseFragment(strings.NewReader(t.text), context)
	return nodes
}

func (t *StringTemplate) MarshalText() (text []byte, err error) {
	return []byte(t.text), nil
}

func (t *StringTemplate) WriteHash(w io.Writer) {
	w.Write([]byte(t.text))
}

func (t *StringTemplate) Min() PreloadableTemplate {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", t.text)
	if err != nil {
		panic(err)
	}

	return PreloadableTemplate{String(text)}
}

func (t *StringTemplate) Preload() Template {
	return Preload(t)
}

func String(text string) *StringTemplate {
	return &StringTemplate{text}
}
