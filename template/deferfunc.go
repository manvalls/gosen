package template

type deferFuncTemplateBuilder struct {
	f func() Template
}

func (t *deferFuncTemplateBuilder) Template() Template {
	return t.f()
}

func DeferFunc(f func() Template) Template {
	return Defer(&deferFuncTemplateBuilder{f})
}
