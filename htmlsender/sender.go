package htmlsender

import (
	"io"
	"sync"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/selectorcache"
	"github.com/manvalls/mutexmap"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type HTMLSender struct {
	mutex         *mutexmap.MutexMap[uint64]
	txMux         *sync.RWMutex
	document      *html.Node
	head          *html.Node
	body          *html.Node
	selectorCache *selectorcache.SelectorCache
	once          map[string]bool
	onceMutex     *sync.Mutex
	idMap         map[string][]*html.Node
	connected     map[*html.Node]bool
}

func NewHTMLSender(cache *selectorcache.SelectorCache) *HTMLSender {
	connected := map[*html.Node]bool{}

	document := &html.Node{
		Type: html.DocumentNode,
	}

	doctype := &html.Node{
		Type: html.DoctypeNode,
		Data: "html",
	}

	document.AppendChild(doctype)

	htmlNode := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Html,
		Data:     "html",
	}

	document.AppendChild(htmlNode)

	head := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Head,
		Data:     "head",
	}

	htmlNode.AppendChild(head)

	body := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}

	htmlNode.AppendChild(body)

	connected[document] = true
	connected[doctype] = true
	connected[htmlNode] = true
	connected[head] = true
	connected[body] = true

	return &HTMLSender{
		mutex:         mutexmap.NewMutexMap[uint64](),
		txMux:         &sync.RWMutex{},
		document:      document,
		head:          head,
		body:          body,
		selectorCache: cache,
		once:          map[string]bool{},
		onceMutex:     &sync.Mutex{},
		idMap:         map[string][]*html.Node{},
		connected:     connected,
	}
}

func (s *HTMLSender) SendCommand(command any) {
	switch c := command.(type) {

	case commands.StartRoutineCommand:
		s.mutex.Lock(c.StartRoutine)
		defer s.mutex.Unlock(c.StartRoutine)
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)

	case commands.TransactionCommand:
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)
		s.transaction(c)

	}
}

func (s *HTMLSender) PrependScript(hydrationData string) {
	script := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Script,
		Data:     "script",
	}

	script.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: hydrationData,
	})

	s.head.InsertBefore(script, s.head.FirstChild)
}

func (s *HTMLSender) Render(w io.Writer) error {
	return html.Render(w, s.document)
}

func getId(node *html.Node) string {
	if node.Type != html.ElementNode {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}

	return ""
}

func (s *HTMLSender) processAddedNode(node *html.Node) {
	s.connected[node] = true

	if id := getId(node); id != "" {
		s.idMap[id] = append(s.idMap[id], node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		s.processAddedNode(child)
	}
}

func (s *HTMLSender) processRemovedNode(node *html.Node) {
	s.connected[node] = false

	if id := getId(node); id != "" {
		for i, n := range s.idMap[id] {
			if n == node {
				s.idMap[id] = append(s.idMap[id][:i], s.idMap[id][i+1:]...)
				break
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		s.processRemovedNode(child)
	}
}

func (s *HTMLSender) processAttrToBeChanged(node *html.Node, attr string, newValue string, oldValue string) {
	if attr != "id" {
		return
	}

	if oldValue != "" {
		for i, n := range s.idMap[oldValue] {
			if n == node {
				s.idMap[oldValue] = append(s.idMap[oldValue][:i], s.idMap[oldValue][i+1:]...)
				break
			}
		}
	}

	if newValue != "" {
		s.idMap[newValue] = append(s.idMap[newValue], node)
	}
}
