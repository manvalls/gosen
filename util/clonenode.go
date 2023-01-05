package util

import (
	"golang.org/x/net/html"
)

func CloneNode(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}

	newNode := &html.Node{
		Type:      node.Type,
		DataAtom:  node.DataAtom,
		Data:      node.Data,
		Namespace: node.Namespace,
	}

	newNode.Attr = make([]html.Attribute, len(node.Attr))
	copy(newNode.Attr, node.Attr)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		newNode.AppendChild(CloneNode(c))
	}

	return newNode
}
