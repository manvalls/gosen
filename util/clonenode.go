package util

import "golang.org/x/net/html"

func CloneNode(node *html.Node, cache map[*html.Node]*html.Node) *html.Node {
	if node == nil {
		return nil
	}

	if val, ok := cache[node]; ok {
		return val
	}

	newNode := &html.Node{}
	cache[node] = newNode

	newNode.Parent = CloneNode(node.Parent, cache)
	newNode.FirstChild = CloneNode(node.FirstChild, cache)
	newNode.LastChild = CloneNode(node.LastChild, cache)
	newNode.PrevSibling = CloneNode(node.PrevSibling, cache)
	newNode.NextSibling = CloneNode(node.NextSibling, cache)

	newNode.Type = node.Type
	newNode.DataAtom = node.DataAtom
	newNode.Data = node.Data
	newNode.Namespace = node.Namespace

	newNode.Attr = make([]html.Attribute, len(node.Attr))
	copy(newNode.Attr, node.Attr)

	return newNode
}
