package template

import (
	"io"

	"golang.org/x/net/html"
)

type ReadCloserFactory interface {
	ReadCloser() (io.ReadCloser, error)
}

type readTemplate struct {
	builder ReadCloserFactory
}

func (t *readTemplate) GetFragment(context *html.Node) []*html.Node {
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		return nil
	}

	fragment, _ := html.ParseFragment(readCloser, context)
	readCloser.Close()
	return fragment
}

func (t *readTemplate) MarshalText() (text []byte, err error) {
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		return nil, err
	}

	text, err = io.ReadAll(readCloser)
	readCloser.Close()
	return text, err
}

func Read(builder ReadCloserFactory) Template {
	return &readTemplate{builder}
}
