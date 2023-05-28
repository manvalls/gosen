package template

import (
	"io"

	"golang.org/x/net/html"
)

type RawTemplate struct {
	text string
}

func (t *RawTemplate) Fragment(context *html.Node) []*html.Node {
	return []*html.Node{
		{
			Type: html.RawNode,
			Data: t.text,
		},
	}
}

func (t *RawTemplate) MarshalText() (text []byte, err error) {
	return []byte(t.text), nil
}

func (t *RawTemplate) WriteHash(w io.Writer) {
	w.Write([]byte(t.text))
}

func (t *RawTemplate) Min() Template {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", t.text)
	if err != nil {
		panic(err)
	}

	return &RawTemplate{text}
}

func Raw(text string) *RawTemplate {
	return &RawTemplate{text}
}
