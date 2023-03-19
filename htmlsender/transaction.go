package htmlsender

// NOTE: There's a lot of duplication between node and content, fix it later.
// This file needs splitting up into multiple files, one per command.

import (
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/template"
	"github.com/manvalls/gosen/util"
	"golang.org/x/net/html"
)

type fragment struct {
	nodes []*html.Node
}

type content struct {
	parent *html.Node
}

func queryFirstFromNode(node *html.Node, selector cascadia.Sel) *html.Node {
	return cascadia.Query(node, selector)
}

func queryFirstFromFragment(f *fragment, selector cascadia.Sel) *html.Node {
	for _, node := range f.nodes {
		if selector.Match(node) {
			return node
		}

		result := queryFirstFromNode(node, selector)
		if result != nil {
			return result
		}
	}

	return nil
}

func queryFirst(parent []any, selector cascadia.Sel) []any {
	nodes := make([]any, 0)

	for _, node := range parent {
		switch n := node.(type) {

		case *html.Node:
			node := queryFirstFromNode(n, selector)
			if node != nil {
				return []any{node}
			}

		case *fragment:
			node := queryFirstFromFragment(n, selector)
			if node != nil {
				return []any{node}
			}

		case content:
			node := queryFirstFromNode(n.parent, selector)
			if node != nil {
				return []any{node}
			}

		}
	}

	return nodes
}

func queryAllFromNode(node *html.Node, selector cascadia.Sel) []*html.Node {
	return cascadia.QueryAll(node, selector)
}

func queryAllFromFragment(f *fragment, selector cascadia.Sel) []*html.Node {
	nodes := make([]*html.Node, 0)

	for _, node := range f.nodes {
		if selector.Match(node) {
			nodes = append(nodes, node)
		}

		result := queryAllFromNode(node, selector)
		nodes = append(nodes, result...)
	}

	return nodes
}

func queryAll(parent []any, selector cascadia.Sel) []any {
	nodes := []*html.Node{}

	for _, node := range parent {
		switch n := node.(type) {

		case *html.Node:
			result := queryAllFromNode(n, selector)
			nodes = append(nodes, result...)

		case *fragment:
			result := queryAllFromFragment(n, selector)
			nodes = append(nodes, result...)

		case content:
			result := queryAllFromNode(n.parent, selector)
			nodes = append(nodes, result...)

		}

	}

	resultArray := make([]any, len(nodes))
	for i, node := range nodes {
		resultArray[i] = node
	}

	return resultArray
}

func getNodesToInsert(nodes map[uint64][]any, id uint64, clone bool) []*html.Node {
	result := []*html.Node{}
	for _, node := range nodes[id] {
		switch n := node.(type) {

		case *html.Node:
			if clone {
				result = append(result, util.CloneNode(n))
			} else {
				if n.Parent != nil {
					n.Parent.RemoveChild(n)
				}

				result = append(result, n)
			}

		case *fragment:
			if clone {
				for _, node := range n.nodes {
					result = append(result, util.CloneNode(node))
				}
			} else {
				result = append(result, n.nodes...)
				n.nodes = nil
			}

		case content:
			if clone {
				for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
					result = append(result, util.CloneNode(c))
				}
			} else {
				for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
					result = append(result, c)
				}

				for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
					n.parent.RemoveChild(c)
				}
			}

		}
	}

	return result
}

func (s *HTMLSender) transaction(c commands.TransactionCommand) {
	if c.Once {
		s.onceMutex.Lock()

		if s.once[c.Hash] {
			s.onceMutex.Unlock()
			return
		}

		s.once[c.Hash] = true
		s.onceMutex.Unlock()
	}

	nodes := make(map[uint64][]any)
	nodes[0] = []any{s.document}

	for _, command := range c.Transaction {

		switch cmd := command.(type) {

		case commands.SelectorSubCommand:

			sel, err := s.selectorCache.Get(cmd.Selector, cmd.Dynamic)
			if err != nil {
				continue
			}

			s.txMux.RLock()
			nodes[cmd.Id] = queryFirst(nodes[cmd.Parent], sel)
			s.txMux.RUnlock()

		case commands.SelectorAllSubCommand:

			sel, err := s.selectorCache.Get(cmd.SelectorAll, cmd.Dynamic)
			if err != nil {
				continue
			}

			s.txMux.RLock()
			nodes[cmd.Id] = queryAll(nodes[cmd.Parent], sel)
			s.txMux.RUnlock()

		case commands.FragmentSubCommand:

			nodes[cmd.Id] = []any{&fragment{template.WithFallback(cmd.Fragment).GetFragment(nil)}}

		case commands.ContentSubCommand:
			result := []any{}

			for _, node := range nodes[cmd.Content] {
				switch n := node.(type) {

				case *html.Node:
					result = append(result, content{n})

				case *fragment:
					result = append(result, n)

				case content:
					result = append(result, n)

				}
			}

			nodes[cmd.Id] = result

		case commands.ParentSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					if n.Parent != nil {
						result = append(result, n.Parent)
					}

				case content:
					result = append(result, n.parent)

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.Parent] = result

		case commands.FirstChildSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					if n.FirstChild != nil {
						result = append(result, n.FirstChild)
					}

				case content:
					if n.parent.FirstChild != nil {
						result = append(result, n.parent.FirstChild)
					}

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.FirstChild] = result

		case commands.LastChildSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					if n.LastChild != nil {
						result = append(result, n.LastChild)
					}

				case content:
					if n.parent.LastChild != nil {
						result = append(result, n.parent.LastChild)
					}

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.LastChild] = result

		case commands.NextSiblingSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					if n.NextSibling != nil {
						result = append(result, n.NextSibling)
					}

				case content:
					if n.parent.NextSibling != nil {
						result = append(result, n.parent.NextSibling)
					}

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.NextSibling] = result

		case commands.PrevSiblingSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					if n.PrevSibling != nil {
						result = append(result, n.PrevSibling)
					}

				case content:
					if n.parent.PrevSibling != nil {
						result = append(result, n.parent.PrevSibling)
					}

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.PrevSibling] = result

		case commands.CloneSubCommand:
			result := []any{}
			s.txMux.RLock()

			for _, node := range nodes[cmd.Clone] {
				switch n := node.(type) {

				case *html.Node:
					result = append(result, util.CloneNode(n))

				case *fragment:
					f := &fragment{}
					for _, node := range n.nodes {
						f.nodes = append(f.nodes, util.CloneNode(node))
					}

					result = append(result, f)

				case content:
					f := &fragment{}
					for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
						f.nodes = append(f.nodes, util.CloneNode(c))
					}

					result = append(result, f)

				}
			}

			s.txMux.RUnlock()
			nodes[cmd.Id] = result

		case commands.TextSubCommand:
			s.txMux.Lock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						n.RemoveChild(c)
					}

					n.AppendChild(&html.Node{
						Type: html.TextNode,
						Data: cmd.Text,
					})

				case *fragment:
					n.nodes = []*html.Node{
						{
							Type: html.TextNode,
							Data: cmd.Text,
						},
					}

				case content:
					for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
						n.parent.RemoveChild(c)
					}

					n.parent.AppendChild(&html.Node{
						Type: html.TextNode,
						Data: cmd.Text,
					})

				}
			}

			s.txMux.Unlock()

		case commands.HtmlSubCommand:

			s.txMux.Lock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						n.RemoveChild(c)
					}

					for _, child := range template.WithFallback(cmd.Html).GetFragment(n) {
						n.AppendChild(child)
					}

				case *fragment:
					n.nodes = template.WithFallback(cmd.Html).GetFragment(nil)

				case content:
					for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
						n.parent.RemoveChild(c)
					}

					for _, child := range template.WithFallback(cmd.Html).GetFragment(n.parent) {
						n.parent.AppendChild(child)
					}

				}
			}

			s.txMux.Unlock()

		case commands.AttrSubCommand:

			s.txMux.Lock()

		loop:
			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for i, attr := range n.Attr {
						if attr.Key == cmd.Attr {
							n.Attr[i].Val = cmd.Value
							continue loop
						}
					}

					n.Attr = append(n.Attr, html.Attribute{
						Key: cmd.Attr,
						Val: cmd.Value,
					})

				}
			}

			s.txMux.Unlock()

		case commands.RemoveAttrSubCommand:
			s.txMux.Lock()

			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for i, attr := range n.Attr {
						if attr.Key == cmd.RemoveAttr {
							n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
							break
						}
					}

				}
			}

			s.txMux.Unlock()

		case commands.AddToAttrSubCommand:
			s.txMux.Lock()

		addToAttrLoop:
			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for i, attr := range n.Attr {
						if attr.Key == cmd.AddToAttr {
							fields := strings.Fields(attr.Val)
							for _, field := range fields {
								if field == cmd.Value {
									continue addToAttrLoop
								}
							}
							n.Attr[i].Val = attr.Val + " " + cmd.Value
							continue addToAttrLoop
						}
					}

					n.Attr = append(n.Attr, html.Attribute{
						Key: cmd.AddToAttr,
						Val: cmd.Value,
					})
				}

			}

			s.txMux.Unlock()

		case commands.RemoveFromAttrSubCommand:
			s.txMux.Lock()

		rmFromAttrLoop:
			for _, node := range nodes[cmd.Target] {
				switch n := node.(type) {

				case *html.Node:
					for i, attr := range n.Attr {
						if attr.Key == cmd.RemoveFromAttr {
							fields := strings.Fields(attr.Val)
							for j, field := range fields {
								if field == cmd.Value {
									fields = append(fields[:j], fields[j+1:]...)
									n.Attr[i].Val = strings.Join(fields, " ")
									continue rmFromAttrLoop
								}
							}
						}
					}

				}
			}

			s.txMux.Unlock()

		case commands.RemoveSubCommand:
			s.txMux.Lock()

			for _, node := range nodes[cmd.Remove] {
				switch n := node.(type) {

				case *html.Node:
					if n.Parent != nil {
						n.Parent.RemoveChild(n)
					}

				case content:
					for c := n.parent.FirstChild; c != nil; c = c.NextSibling {
						n.parent.RemoveChild(c)
					}

				}
			}

			s.txMux.Unlock()

		case commands.InsertNodeBeforeSubCommand:
			s.txMux.Lock()

			parents := nodes[cmd.Parent]
		inbloop:
			for _, node := range parents {
				switch n := node.(type) {

				case *html.Node:
					nodesToInsert := getNodesToInsert(nodes, cmd.InsertNodeBefore, len(parents) > 1)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok || rn.Parent != n {
							continue
						}

						for _, ni := range nodesToInsert {
							n.InsertBefore(ni, rn)
						}

						continue inbloop
					}

					for _, ni := range nodesToInsert {
						n.AppendChild(ni)
					}

				case content:
					nodesToInsert := getNodesToInsert(nodes, cmd.InsertNodeBefore, len(parents) > 1)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok || rn.Parent != n.parent {
							continue
						}

						for _, ni := range nodesToInsert {
							n.parent.InsertBefore(ni, rn)
						}

						continue inbloop
					}

					for _, ni := range nodesToInsert {
						n.parent.AppendChild(ni)
					}

				case *fragment:
					nodesToInsert := getNodesToInsert(nodes, cmd.InsertNodeBefore, len(parents) > 1)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok {
							continue
						}

						var i int
						var c *html.Node
						var found bool
						for i, c = range n.nodes {
							if rn == c {
								found = true
								break
							}
						}

						if !found {
							continue
						}

						n.nodes = append(n.nodes[:i], append(nodesToInsert, n.nodes[i:]...)...)
						continue inbloop
					}

					n.nodes = append(n.nodes, nodesToInsert...)

				}
			}

			s.txMux.Unlock()

		case commands.InsertBeforeSubCommand:
			s.txMux.Lock()

			parents := nodes[cmd.Parent]
		ploop:
			for _, node := range parents {
				switch n := node.(type) {

				case *html.Node:
					nodesToInsert := cmd.InsertBefore.GetFragment(n)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok || rn.Parent != n {
							continue
						}

						for _, ni := range nodesToInsert {
							n.InsertBefore(ni, rn)
						}

						continue ploop
					}

					for _, ni := range nodesToInsert {
						n.AppendChild(ni)
					}

				case content:
					nodesToInsert := cmd.InsertBefore.GetFragment(n.parent)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok || rn.Parent != n.parent {
							continue
						}

						for _, ni := range nodesToInsert {
							n.parent.InsertBefore(ni, rn)
						}

						continue ploop
					}

					for _, ni := range nodesToInsert {
						n.parent.AppendChild(ni)
					}

				case *fragment:
					nodesToInsert := cmd.InsertBefore.GetFragment(nil)

					for _, r := range nodes[cmd.Ref] {
						rn, ok := r.(*html.Node)
						if !ok {
							continue
						}

						var i int
						var c *html.Node
						var found bool
						for i, c = range n.nodes {
							if rn == c {
								found = true
								break
							}
						}

						if !found {
							continue
						}

						n.nodes = append(n.nodes[:i], append(nodesToInsert, n.nodes[i:]...)...)
						continue ploop
					}

					n.nodes = append(n.nodes, nodesToInsert...)

				}
			}

			s.txMux.Unlock()

		case commands.AppendNodeSubCommand:
			s.txMux.Lock()

			parents := nodes[cmd.Parent]
			for _, node := range parents {
				switch n := node.(type) {

				case *html.Node:
					nodesToInsert := getNodesToInsert(nodes, cmd.AppendNode, len(parents) > 1)
					for _, ni := range nodesToInsert {
						n.AppendChild(ni)
					}

				case content:
					nodesToInsert := getNodesToInsert(nodes, cmd.AppendNode, len(parents) > 1)
					for _, ni := range nodesToInsert {
						n.parent.AppendChild(ni)
					}

				case *fragment:
					nodesToInsert := getNodesToInsert(nodes, cmd.AppendNode, len(parents) > 1)
					n.nodes = append(n.nodes, nodesToInsert...)

				}
			}

			s.txMux.Unlock()

		case commands.AppendSubCommand:
			s.txMux.Lock()

			parents := nodes[cmd.Parent]
			for _, node := range parents {
				switch n := node.(type) {

				case *html.Node:
					nodesToInsert := cmd.Append.GetFragment(n)
					for _, ni := range nodesToInsert {
						n.AppendChild(ni)
					}

				case content:
					nodesToInsert := cmd.Append.GetFragment(n.parent)
					for _, ni := range nodesToInsert {
						n.parent.AppendChild(ni)
					}

				case *fragment:
					nodesToInsert := cmd.Append.GetFragment(nil)
					n.nodes = append(n.nodes, nodesToInsert...)

				}
			}

			s.txMux.Unlock()

		}

	}
}
