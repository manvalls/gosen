package template

import "golang.org/x/net/html"

type Empty struct{}

func (e Empty) GetFragment(context *html.Node) []*html.Node {
	return nil
}

func (e Empty) MarshalText() (text []byte, err error) {
	return nil, nil
}
