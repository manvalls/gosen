package template

import "io"

type funcReadCloserFactory struct {
	f func() (io.ReadCloser, error)
}

func (f *funcReadCloserFactory) ReadCloser() (io.ReadCloser, error) {
	return f.f()
}

func ReadFunc(f func() (io.ReadCloser, error)) Template {
	return Read(&funcReadCloserFactory{f})
}
