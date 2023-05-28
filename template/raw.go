package template

import (
	"html/template"
	"io"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type RawTemplate struct {
	text string
	mux  *sync.Mutex
	tpl  *template.Template
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

func (t *RawTemplate) Min() *RawTemplate {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", t.text)
	if err != nil {
		panic(err)
	}

	return &RawTemplate{text, &sync.Mutex{}, nil}
}

func (t *RawTemplate) getTpl() *template.Template {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.tpl == nil {
		t.tpl = template.Must(template.New("").Parse(t.text))
	}
	return t.tpl
}

func (t *RawTemplate) Execute(data any) Template {
	b := &strings.Builder{}
	t.getTpl().Execute(b, data)
	return Raw(b.String())
}

func Raw(text string) *RawTemplate {
	return &RawTemplate{text, &sync.Mutex{}, nil}
}
