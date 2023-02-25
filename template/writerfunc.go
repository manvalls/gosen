package template

import "io"

type writerFuncWriterTo struct {
	f func(io.Writer) error
}

func (w *writerFuncWriterTo) WriteTo(writer io.Writer) (int64, error) {
	return 0, w.f(writer)
}

func WriterFunc(f func(io.Writer) error) Template {
	return WriterTo(&writerFuncWriterTo{f})
}
