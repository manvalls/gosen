package template

func WithFallback(template Template) Template {
	if template == nil {
		return Empty{}
	}

	return template
}
