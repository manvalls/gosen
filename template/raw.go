package template

import "golang.org/x/net/html"

type RawTemplate struct {
	text string
}

func (t *RawTemplate) GetFragment(context *html.Node) []*html.Node {
	return []*html.Node{
		{
			Type: html.RawNode,
			Data: t.text,
		},
	}
}

func (t *RawTemplate) MarshalText() (text []byte, err error) {
	return []byte(t.text), nil
}

func (t *RawTemplate) Min() Template {
	minifierMutex.Lock()
	defer minifierMutex.Unlock()

	m := getMinifier()
	text, err := m.String("text/html", t.text)
	if err != nil {
		panic(err)
	}

	return &RawTemplate{text}
}

func Raw(text string) Template {
	return &RawTemplate{text}
}
