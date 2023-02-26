package template

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type WriterToTemplate struct {
	writerTo io.WriterTo
}

func (t *WriterToTemplate) GetFragment(context *html.Node) []*html.Node {
	reader, writer := io.Pipe()
	go func() {
		t.writerTo.WriteTo(writer)
		writer.Close()
	}()

	fragment, _ := html.ParseFragment(reader, context)
	return fragment
}

func (t *WriterToTemplate) MarshalText() (text []byte, err error) {
	buffer := &bytes.Buffer{}
	_, err = t.writerTo.WriteTo(buffer)
	return buffer.Bytes(), err
}

func (t *WriterToTemplate) WriteHash(w io.Writer) {
	t.writerTo.WriteTo(w)
}

func (t *WriterToTemplate) Min() PreloadableTemplate {
	return PreloadableTemplate{&WriterToTemplate{&minWriterTo{t.writerTo}}}
}

func (t *WriterToTemplate) Preload() Template {
	return Preload(t)
}

func WriterTo(writerTo io.WriterTo) *WriterToTemplate {
	return &WriterToTemplate{writerTo}
}

type minWriterTo struct {
	io.WriterTo
}

func (mw *minWriterTo) WriteTo(w io.Writer) (int64, error) {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	mfw := m.Writer("text/html", w)

	n, err := mw.WriterTo.WriteTo(mfw)
	mfw.Close()
	return n, err
}
