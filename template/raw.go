package template

import "golang.org/x/net/html"

type rawTemplate struct {
	text string
}

func (t *rawTemplate) GetFragment(context *html.Node) []*html.Node {
	return []*html.Node{
		{
			Type: html.RawNode,
			Data: t.text,
		},
	}
}

func (t *rawTemplate) MarshalText() (text []byte, err error) {
	return []byte(t.text), nil
}

func Raw(text string) Template {
	return &rawTemplate{text}
}
