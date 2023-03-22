package template

import (
	"encoding/binary"
	"io"

	"github.com/manvalls/gosen/util"
	"github.com/zeebo/xxh3"
	"golang.org/x/net/html"
)

type preloadedTemplate struct {
	fragment []*html.Node
	text     []byte
	hash     []byte
}

func (t *preloadedTemplate) Fragment(context *html.Node) []*html.Node {
	result := make([]*html.Node, len(t.fragment))
	for i, node := range t.fragment {
		result[i] = util.CloneNode(node)
	}

	return result
}

func (t *preloadedTemplate) MarshalText() (text []byte, err error) {
	return t.text, nil
}

func (t *preloadedTemplate) WriteHash(w io.Writer) {
	w.Write(t.hash)
}

func Preload(t Template) Template {
	fragment := t.Fragment(nil)
	text, err := t.MarshalText()
	if err != nil {
		panic(err)
	}

	h := xxh3.Hash(text)
	hash := make([]byte, 8)
	binary.PutUvarint(hash, h)

	return &preloadedTemplate{
		fragment: fragment,
		text:     text,
		hash:     hash,
	}
}

type PreloadableTemplate struct {
	Template
}

func (t PreloadableTemplate) Preload() Template {
	return Preload(t.Template)
}
