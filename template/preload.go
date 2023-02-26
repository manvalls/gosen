package template

import (
	"github.com/manvalls/gosen/util"
	"golang.org/x/net/html"
)

type preloadedTemplate struct {
	fragment []*html.Node
	text     []byte
}

func (t *preloadedTemplate) GetFragment(context *html.Node) []*html.Node {
	result := make([]*html.Node, len(t.fragment))
	for i, node := range t.fragment {
		result[i] = util.CloneNode(node)
	}

	return result
}

func (t *preloadedTemplate) MarshalText() (text []byte, err error) {
	return t.text, nil
}

func Preload(t Template) Template {
	fragment := t.GetFragment(nil)
	text, err := t.MarshalText()
	if err != nil {
		panic(err)
	}

	return &preloadedTemplate{
		fragment: fragment,
		text:     text,
	}
}

type PreloadableTemplate struct {
	Template
}

func (t PreloadableTemplate) Preload() Template {
	return Preload(t.Template)
}
