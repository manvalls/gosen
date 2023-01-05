package template

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

type writeTemplate struct {
	writerTo io.WriterTo
}

func (t *writeTemplate) GetFragment(context *html.Node) []*html.Node {
	reader, writer := io.Pipe()
	go func() {
		t.writerTo.WriteTo(writer)
		writer.Close()
	}()

	fragment, _ := html.ParseFragment(reader, context)
	return fragment
}

func (t *writeTemplate) MarshalText() (text []byte, err error) {
	buffer := &bytes.Buffer{}
	_, err = t.writerTo.WriteTo(buffer)
	return buffer.Bytes(), err
}

func Write(writerTo io.WriterTo) Template {
	return &writeTemplate{writerTo}
}
