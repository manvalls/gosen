package commands

import "github.com/manvalls/gosen/template"

type Node struct {
	tx *Transaction
	id uint
}

// Selectors

func (n Node) S(selector string, args ...interface{}) Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(SelectorSubCommand{nextId, selector, args, n.id})
	return Node{n.tx, nextId}
}

func (n Node) All(selector string, args ...interface{}) Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(SelectorAllSubCommand{nextId, selector, args, n.id})
	return Node{n.tx, nextId}
}

type ContentSubCommand struct {
	Id      uint `json:"id"`
	Content uint `json:"content"`
}

func (n Node) Content() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(ContentSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type ParentSubCommand struct {
	Parent uint `json:"parent"`
	Target uint `json:"target"`
}

func (n Node) Parent() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(ParentSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type FirstChildSubCommand struct {
	FirstChild uint `json:"firstChild"`
	Target     uint `json:"target"`
}

func (n Node) FirstChild() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(FirstChildSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type LastChildSubCommand struct {
	LastChild uint `json:"lastChild"`
	Target    uint `json:"target"`
}

func (n Node) LastChild() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(LastChildSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type NextSiblingSubCommand struct {
	NextSibling uint `json:"nextSibling"`
	Target      uint `json:"target"`
}

func (n Node) NextSibling() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(NextSiblingSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type PrevSiblingSubCommand struct {
	PrevSibling uint `json:"prevSibling"`
	Target      uint `json:"target"`
}

func (n Node) PrevSibling() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(PrevSiblingSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

// Creation

type CloneSubCommand struct {
	Id    uint `json:"id"`
	Clone uint `json:"clone"`
}

func (n Node) Clone() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(CloneSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

// InnerText and InnerHTML

type TextSubCommand struct {
	Target uint   `json:"target"`
	Text   string `json:"text"`
}

func (n Node) Text(text string) Node {
	n.tx.SendCommand(TextSubCommand{n.id, text})
	return n
}

type HtmlSubCommand struct {
	Target uint              `json:"target"`
	Html   template.Template `json:"html"`
}

func (n Node) Html(html template.Template) Node {
	n.tx.SendCommand(HtmlSubCommand{n.id, html})
	return n
}

// Attributes

type AttrSubCommand struct {
	Target uint   `json:"target"`
	Attr   string `json:"attr"`
	Value  string `json:"value"`
}

func (n Node) Attr(name string, value string) Node {
	n.tx.SendCommand(AttrSubCommand{n.id, name, value})
	return n
}

type RmAttrSubCommand struct {
	Target uint   `json:"target"`
	RmAttr string `json:"rmAttr"`
}

func (n Node) RmAttr(name string) Node {
	n.tx.SendCommand(RmAttrSubCommand{n.id, name})
	return n
}

// Node manipulation

type RemoveSubCommand struct {
	Remove uint `json:"remove"`
}

func (n Node) Remove() {
	n.tx.SendCommand(RemoveSubCommand{n.id})
}

type InsertBeforeSubCommand struct {
	Parent       uint `json:"parent"`
	Ref          uint `json:"ref"`
	InsertBefore any  `json:"insertBefore"`
}

func (n Node) InsertBefore(child any, ref Node) Node {
	switch c := child.(type) {

	case Node:
		if n.tx != c.tx || n.tx != ref.tx {
			panic("Nodes must belong to the same transaction")
		}

		n.tx.SendCommand(InsertBeforeSubCommand{n.id, ref.id, c.id})

	case template.Template:
		n.tx.SendCommand(InsertBeforeSubCommand{n.id, ref.id, c})

	}

	return n
}

type AppendSubCommand struct {
	Parent uint `json:"parent"`
	Append any  `json:"append"`
}

func (n Node) Append(child any) Node {
	switch c := child.(type) {

	case Node:
		if n.tx != c.tx {
			panic("Nodes must belong to the same transaction")
		}

		n.tx.SendCommand(AppendSubCommand{n.id, c.id})

	case template.Template:
		n.tx.SendCommand(AppendSubCommand{n.id, c})

	}

	return n
}

// Misc

type WaitSubCommand struct {
	Target  uint   `json:"target"`
	Wait    string `json:"wait"`
	Timeout uint   `json:"timeout,omitempty"`
}

func (n Node) Wait(event string, timeout uint) Node {
	n.tx.SendCommand(WaitSubCommand{n.id, event, timeout})
	return n
}
