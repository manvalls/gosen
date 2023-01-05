package template

type deferFuncTemplateFactory struct {
	f func() Template
}

func (t *deferFuncTemplateFactory) Template() Template {
	return t.f()
}

func DeferFunc(f func() Template) Template {
	return Defer(&deferFuncTemplateFactory{f})
}
