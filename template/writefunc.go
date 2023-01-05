package template

import "io"

type writeFuncWriterTo struct {
	f func(io.Writer) error
}

func (w *writeFuncWriterTo) WriteTo(writer io.Writer) (int64, error) {
	return 0, w.f(writer)
}

func WriteFunc(f func(io.Writer) error) Template {
	return Write(&writeFuncWriterTo{f})
}
