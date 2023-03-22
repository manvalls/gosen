package template

import (
	"io"

	"golang.org/x/net/html"
)

type Empty struct{}

func (e Empty) Fragment(context *html.Node) []*html.Node {
	return nil
}

func (e Empty) MarshalText() (text []byte, err error) {
	return nil, nil
}

func (e Empty) WriteHash(w io.Writer) {}
