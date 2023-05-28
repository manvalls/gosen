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
type Event = commands.Event

type Template = template.Template
type TemplateFactory = template.TemplateFactory
type ReadCloserFactory = template.ReadCloserFactory
type Empty = template.Empty

func Defer(t TemplateFactory) Template {
	return template.Defer(t)
}

func DeferFunc(f func() Template) Template {
	return template.DeferFunc(f)
}

func File(fs fs.FS, name string) *template.ReadTemplate {
	return template.File(fs, name)
}

func Raw(s string) *template.RawTemplate {
	return template.Raw(s)
}

func RawFile(fs fs.FS, name string) *template.RawFileTemplate {
	return template.RawFile(fs, name)
}

func String(s string) *template.StringTemplate {
	return template.String(s)
}

func Read(r ReadCloserFactory) *template.ReadTemplate {
	return template.Read(r)
}

func ReadFunc(f func() (io.ReadCloser, error)) *template.ReadTemplate {
	return template.ReadFunc(f)
}

func WithFallback(t Template) Template {
	return template.WithFallback(t)
}

func WriterTo(w io.WriterTo) *template.WriterToTemplate {
	return template.WriterTo(w)
}

func WriterFunc(f func(io.Writer) error) *template.WriterToTemplate {
	return template.WriterFunc(f)
}

func Preload(t Template) Template {
	return template.Preload(t)
}
