package template

import (
	"io"
	"io/fs"
)

type fileReadCloserFactory struct {
	fs       fs.FS
	fileName string
}

func (f *fileReadCloserFactory) ReadCloser() (io.ReadCloser, error) {
	return f.fs.Open(f.fileName)
}

func File(fs fs.FS, fileName string) *ReadTemplate {
	return Read(&fileReadCloserFactory{fs, fileName})
}
