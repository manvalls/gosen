package template

import (
	"sync"

	"github.com/manvalls/gosen/util"
	"golang.org/x/net/html"
)

type cachedTemplate struct {
	template    Template
	fragmentMux sync.Mutex
	fragment    []*html.Node
	textMux     sync.Mutex
	text        []byte
}

func (t *cachedTemplate) GetFragment(context *html.Node) []*html.Node {
	t.fragmentMux.Lock()
	defer t.fragmentMux.Unlock()
	if t.fragment == nil {
		t.fragment = t.template.GetFragment(context)
	}

	fragment := make([]*html.Node, len(t.fragment))
	for i, node := range t.fragment {
		fragment[i] = util.CloneNode(node)
	}

	return fragment
}

func (t *cachedTemplate) MarshalText() (text []byte, err error) {
	t.textMux.Lock()
	defer t.textMux.Unlock()

	if t.text == nil {
		t.text, err = t.template.MarshalText()
	}

	return t.text, err
}

func Cache(template Template) Template {
	return &cachedTemplate{
		template: template,
	}
}
