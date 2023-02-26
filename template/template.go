package template

import (
	"io"

	"golang.org/x/net/html"
)

type Template interface {
	GetFragment(context *html.Node) []*html.Node
	MarshalText() (text []byte, err error)
	WriteHash(w io.Writer)
}
