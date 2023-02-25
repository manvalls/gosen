package template

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type writerToTemplate struct {
	writerTo io.WriterTo
}

func (t *writerToTemplate) GetFragment(context *html.Node) []*html.Node {
	reader, writer := io.Pipe()
	go func() {
		t.writerTo.WriteTo(writer)
		writer.Close()
	}()

	fragment, _ := html.ParseFragment(reader, context)
	return fragment
}

func (t *writerToTemplate) MarshalText() (text []byte, err error) {
	buffer := &bytes.Buffer{}
	_, err = t.writerTo.WriteTo(buffer)
	return buffer.Bytes(), err
}

func WriterTo(writerTo io.WriterTo) Template {
	return &writerToTemplate{writerTo}
}
