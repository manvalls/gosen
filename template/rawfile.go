package template

import (
	"io"
	"io/fs"
)

type rawFileTemplateFactory struct {
	fs       fs.FS
	fileName string
}

func (t *rawFileTemplateFactory) Template() Template {
	file, err := t.fs.Open(t.fileName)
	if err != nil {
		return nil
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	return Raw(string(data))
}

func RawFile(fs fs.FS, fileName string) Template {
	return Defer(&rawFileTemplateFactory{fs, fileName})
}
