package util

import (
	"golang.org/x/net/html"
)

func CloneNode(node *html.Node) *html.Node {
	cache := map[*html.Node]*html.Node{}
	cache[node.Parent] = nil
	return cloneNode(node, cache)
}

func cloneNode(node *html.Node, cache map[*html.Node]*html.Node) *html.Node {
	if node == nil {
		return nil
	}

	if val, ok := cache[node]; ok {
		return val
	}

	newNode := &html.Node{}
	cache[node] = newNode

	newNode.Parent = cloneNode(node.Parent, cache)
	newNode.FirstChild = cloneNode(node.FirstChild, cache)
	newNode.LastChild = cloneNode(node.LastChild, cache)
	newNode.PrevSibling = cloneNode(node.PrevSibling, cache)
	newNode.NextSibling = cloneNode(node.NextSibling, cache)

	newNode.Type = node.Type
	newNode.DataAtom = node.DataAtom
	newNode.Data = node.Data
	newNode.Namespace = node.Namespace

	newNode.Attr = make([]html.Attribute, len(node.Attr))
	copy(newNode.Attr, node.Attr)

	return newNode
}
