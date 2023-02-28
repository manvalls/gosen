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
	document      *html.Node
	selectorCache *selectorcache.SelectorCache
	once          map[string]bool
	onceMutex     *sync.Mutex
}

func NewHTMLSender(cache *selectorcache.SelectorCache) *HTMLSender {
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

	return &HTMLSender{
		mutex:         mutexmap.NewMutexMap[uint64](),
		document:      document,
		selectorCache: cache,
		once:          map[string]bool{},
		onceMutex:     &sync.Mutex{},
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

func (s *HTMLSender) Render(w io.Writer) error {
	return html.Render(w, s.document)
}
