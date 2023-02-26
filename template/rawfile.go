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
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return Raw(string(data))
}

type RawFileTemplate struct {
	Template
	fs       fs.FS
	fileName string
}

func (t *RawFileTemplate) Min() *MinRawFileTemplate {
	return &MinRawFileTemplate{
		Template: Defer(&minRawFileTemplateFactory{t.fs, t.fileName}),
		fs:       t.fs,
		fileName: t.fileName,
	}
}

func (t *RawFileTemplate) Preload() Template {
	file, err := t.fs.Open(t.fileName)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return Raw(string(data))
}

func RawFile(fs fs.FS, fileName string) *RawFileTemplate {
	return &RawFileTemplate{
		Template: Defer(&rawFileTemplateFactory{fs, fileName}),
		fs:       fs,
		fileName: fileName,
	}
}

type minRawFileTemplateFactory struct {
	fs       fs.FS
	fileName string
}

func (t *minRawFileTemplateFactory) Template() Template {
	file, err := t.fs.Open(t.fileName)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", string(data))
	if err != nil {
		panic(err)
	}

	return Raw(text)
}

type MinRawFileTemplate struct {
	Template
	fs       fs.FS
	fileName string
}

func (t *MinRawFileTemplate) Preload() Template {
	file, err := t.fs.Open(t.fileName)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", string(data))
	if err != nil {
		panic(err)
	}

	return Raw(text)
}
