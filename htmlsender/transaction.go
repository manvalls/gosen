package htmlsender

import (
	"github.com/andybalholm/cascadia"
	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/util"
	"golang.org/x/net/html"
)

type SelectedNodes struct {
	isFragment bool
	nodes      []*html.Node
}

func queryFirst(parent *SelectedNodes, selector cascadia.Sel) *SelectedNodes {
	if parent.isFragment {
		for _, node := range parent.nodes {
			if selector.Match(node) {
				return &SelectedNodes{
					isFragment: false,
					nodes:      []*html.Node{node},
				}
			}
		}
	}

	for _, node := range parent.nodes {
		result := cascadia.Query(node, selector)
		if result != nil {
			return &SelectedNodes{
				isFragment: false,
				nodes:      []*html.Node{result},
			}
		}
	}

	return &SelectedNodes{
		isFragment: false,
		nodes:      []*html.Node{},
	}
}

func queryAll(parent *SelectedNodes, selector cascadia.Sel) *SelectedNodes {
	nodes := make([]*html.Node, 0)

	if parent.isFragment {
		for _, node := range parent.nodes {
			if selector.Match(node) {
				nodes = append(nodes, node)
			}
		}
	}

	for _, node := range parent.nodes {
		result := cascadia.QueryAll(node, selector)
		nodes = append(nodes, result...)
	}

	return &SelectedNodes{
		isFragment: false,
		nodes:      nodes,
	}
}

func (s *HTMLSender) transaction(c commands.TransactionCommand) {

	nodes := make(map[uint]*SelectedNodes)
	nodes[0] = &SelectedNodes{
		isFragment: false,
		nodes:      []*html.Node{s.document},
	}

	for _, command := range c.Transaction {

		switch cmd := command.(type) {

		case commands.SelectorSubCommand:

			parent := nodes[cmd.Parent]
			if parent == nil {
				continue
			}

			sel, err := s.selectorCache.Get(cmd.Selector, cmd.Args)
			if err != nil {
				continue
			}

			nodes[cmd.Id] = queryFirst(parent, sel)

		case commands.SelectorAllSubCommand:

			parent := nodes[cmd.Parent]
			if parent == nil {
				continue
			}

			sel, err := s.selectorCache.Get(cmd.SelectorAll, cmd.Args)
			if err != nil {
				continue
			}

			nodes[cmd.Id] = queryAll(parent, sel)

		case commands.FragmentSubCommand:

			nodes[cmd.Id] = &SelectedNodes{
				isFragment: true,
				nodes:      cmd.Fragment.GetFragment(nil),
			}

		case commands.ContentSubCommand:
			parent := nodes[cmd.Id]
			if parent == nil {
				continue
			}

			if parent.isFragment {
				nodes[cmd.Content] = &SelectedNodes{
					isFragment: true,
					nodes:      parent.nodes,
				}

				continue
			}

			content := []*html.Node{}
			for _, node := range parent.nodes {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					content = append(content, c)
				}
			}

			nodes[cmd.Content] = &SelectedNodes{
				isFragment: true,
				nodes:      content,
			}

		case commands.CloneSubCommand:
			parent := nodes[cmd.Id]
			if parent == nil {
				continue
			}

			cache := make(map[*html.Node]*html.Node)
			result := make([]*html.Node, len(parent.nodes))
			for i, node := range parent.nodes {
				result[i] = util.CloneNode(node, cache)
			}

			nodes[cmd.Clone] = &SelectedNodes{
				isFragment: parent.isFragment,
				nodes:      result,
			}

		case commands.TextSubCommand:
			parent := nodes[cmd.Target]
			if parent == nil {
				continue
			}

			if parent.isFragment {
				parent.nodes = []*html.Node{
					{
						Type: html.TextNode,
						Data: cmd.Text,
					},
				}
				continue
			}

			for _, node := range parent.nodes {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					node.RemoveChild(c)
				}

				node.AppendChild(&html.Node{
					Type: html.TextNode,
					Data: cmd.Text,
				})
			}

		case commands.HtmlSubCommand:
			parent := nodes[cmd.Target]
			if parent == nil {
				continue
			}

			if parent.isFragment {
				parent.nodes = cmd.Html.GetFragment(nil)
				continue
			}

			for _, node := range parent.nodes {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					node.RemoveChild(c)
				}

				for _, child := range cmd.Html.GetFragment(node) {
					node.AppendChild(child)
				}
			}

		case commands.AttrSubCommand:
			parent := nodes[cmd.Target]
			if parent == nil || parent.isFragment {
				continue
			}

		loop:
			for _, node := range parent.nodes {
				for i, attr := range node.Attr {
					if attr.Key == cmd.Attr {
						node.Attr[i].Val = cmd.Value
						continue loop
					}
				}

				node.Attr = append(node.Attr, html.Attribute{
					Key: cmd.Attr,
					Val: cmd.Value,
				})
			}

		case commands.RmAttrSubCommand:
			parent := nodes[cmd.Target]
			if parent == nil || parent.isFragment {
				continue
			}

			for _, node := range parent.nodes {
				for i, attr := range node.Attr {
					if attr.Key == cmd.RmAttr {
						node.Attr = append(node.Attr[:i], node.Attr[i+1:]...)
						break
					}
				}
			}

		case commands.RemoveSubCommand:
			parent := nodes[cmd.Remove]
			if parent == nil || parent.isFragment {
				continue
			}

			for _, node := range parent.nodes {
				node.Parent.RemoveChild(node)
			}

		case commands.InsertBeforeSubCommand:
			parent := nodes[cmd.Parent]
			ref := nodes[cmd.Ref]
			newChild := nodes[cmd.InsertBefore]

			if parent == nil || newChild == nil {
				continue
			}

			if parent.isFragment {
			top:
				for _, node := range parent.nodes {
					for i, refNode := range ref.nodes {
						if refNode == node {
							for _, new := range newChild.nodes {
								if new.Parent != nil {
									new.Parent.RemoveChild(new)
								}
							}

							parent.nodes = append(parent.nodes[:i], append(newChild.nodes, parent.nodes[i:]...)...)
							break top
						}
					}
				}

				continue
			}

			clone := false
			for _, node := range parent.nodes {
				for _, refNode := range ref.nodes {
					if refNode.Parent != node {
						continue
					}

					for _, new := range newChild.nodes {
						if new.Parent != nil {
							new.Parent.RemoveChild(new)
						}

						if clone {
							new = util.CloneNode(new, make(map[*html.Node]*html.Node))
						}

						node.InsertBefore(new, refNode)
					}

					clone = true
				}
			}

		case commands.AppendSubCommand:
			// TODO

		}

	}

}
