package template

import (
	"io"
	"sync"

	"golang.org/x/net/html"
)

type TemplateFactory interface {
	Template() Template
}

type deferredTemplate struct {
	builder  TemplateFactory
	mux      sync.Mutex
	template Template
}

func (t *deferredTemplate) getTemplate() Template {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.template == nil {
		t.template = WithFallback(t.builder.Template())
	}
	return t.template
}

func (t *deferredTemplate) Fragment(context *html.Node) []*html.Node {
	return t.getTemplate().Fragment(context)
}

func (t *deferredTemplate) MarshalText() (text []byte, err error) {
	return t.getTemplate().MarshalText()
}

func (t *deferredTemplate) WriteHash(w io.Writer) {
	t.getTemplate().WriteHash(w)
}

func Defer(builder TemplateFactory) Template {
	return &deferredTemplate{builder, sync.Mutex{}, nil}
}
