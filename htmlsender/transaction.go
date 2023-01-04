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
			// TODO
		case commands.AttrSubCommand:
			// TODO
		case commands.RmAttrSubCommand:
			// TODO
		case commands.AddToAttrSubCommand:
			// TODO
		case commands.RmFromAttrSubCommand:
			// TODO
		case commands.AddClassSubCommand:
			// TODO
		case commands.RmClassSubCommand:
			// TODO
		case commands.RemoveSubCommand:
			// TODO
		case commands.EmptySubCommand:
			// TODO
		case commands.ReplaceWithSubCommand:
			// TODO
		case commands.InsertBeforeSubCommand:
			// TODO
		case commands.InsertAfterSubCommand:
			// TODO
		case commands.AppendSubCommand:
			// TODO
		case commands.PrependSubCommand:
			// TODO
		}

	}

}
