package gosen

import (
	"io"
	"io/fs"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/template"
)

type Node = commands.Node
type Transaction = commands.Transaction
type Routine = commands.Routine

type Template = template.Template
type TemplateFactory = template.TemplateFactory
type ReadCloserFactory = template.ReadCloserFactory
type Empty = template.Empty

func Cache(t Template) Template {
	return template.Cache(t)
}

func Defer(t TemplateFactory) Template {
	return template.Defer(t)
}

func DeferFunc(f func() Template) Template {
	return template.DeferFunc(f)
}

func File(fs fs.FS, name string) Template {
	return template.File(fs, name)
}

func Raw(s string) Template {
	return template.Raw(s)
}

func RawFile(fs fs.FS, name string) Template {
	return template.RawFile(fs, name)
}

func String(s string) Template {
	return template.String(s)
}

func Read(r ReadCloserFactory) Template {
	return template.Read(r)
}

func ReadFunc(f func() (io.ReadCloser, error)) Template {
	return template.ReadFunc(f)
}

func WithFallback(t Template) Template {
	return template.WithFallback(t)
}

func Write(w io.WriterTo) Template {
	return template.Write(w)
}

func WriteFunc(f func(io.Writer) error) Template {
	return template.WriteFunc(f)
}
