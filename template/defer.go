package template

import (
	"sync"

	"golang.org/x/net/html"
)

type TemplateBuilder interface {
	Template() Template
}

type deferredTemplate struct {
	builder  TemplateBuilder
	mux      sync.Mutex
	template Template
}

func (t *deferredTemplate) getTemplate() Template {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.template == nil {
		t.template = t.builder.Template()
	}
	return t.template
}

func (t *deferredTemplate) GetFragment(context *html.Node) []*html.Node {
	return t.getTemplate().GetFragment(context)
}

func (t *deferredTemplate) MarshalText() (text []byte, err error) {
	return t.getTemplate().MarshalText()
}

func Defer(builder TemplateBuilder) Template {
	return &deferredTemplate{builder, sync.Mutex{}, nil}
}
