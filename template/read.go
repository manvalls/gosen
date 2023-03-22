package template

import (
	"io"
	"io/ioutil"

	"golang.org/x/net/html"
)

type ReadCloserFactory interface {
	ReadCloser() (io.ReadCloser, error)
}

type ReadTemplate struct {
	builder ReadCloserFactory
}

func (t *ReadTemplate) Fragment(context *html.Node) []*html.Node {
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		return nil
	}

	fragment, _ := html.ParseFragment(readCloser, context)
	readCloser.Close()
	return fragment
}

func (t *ReadTemplate) MarshalText() (text []byte, err error) {
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		return nil, err
	}

	text, err = io.ReadAll(readCloser)
	readCloser.Close()
	return text, err
}

func (t *ReadTemplate) WriteHash(w io.Writer) {
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(w, readCloser)
	readCloser.Close()
	if err != nil {
		panic(err)
	}
}

func (t *ReadTemplate) Min() PreloadableTemplate {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	readCloser, err := t.builder.ReadCloser()
	if err != nil {
		panic(err)
	}

	r := m.Reader("text/html", readCloser)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return PreloadableTemplate{String(string(bytes))}
}

func (t *ReadTemplate) Preload() Template {
	return Preload(t)
}

func Read(builder ReadCloserFactory) *ReadTemplate {
	return &ReadTemplate{builder}
}
