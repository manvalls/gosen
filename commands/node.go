package commands

import (
	"fmt"

	"github.com/manvalls/gosen/template"
)

type Node struct {
	tx *Transaction
	id uint64
}

// Selectors

func (n Node) S(selector string, args ...any) Node {

	nextId := n.tx.getNextId()

	if len(args) > 0 {
		n.tx.SendCommand(SelectorSubCommand{nextId, fmt.Sprintf(selector, args...), true, n.id})
		return Node{n.tx, nextId}
	}

	n.tx.SendCommand(SelectorSubCommand{nextId, selector, false, n.id})
	return Node{n.tx, nextId}
}

func (n Node) All(selector string, args ...any) Node {
	nextId := n.tx.getNextId()

	if len(args) > 0 {
		n.tx.SendCommand(SelectorAllSubCommand{nextId, fmt.Sprintf(selector, args...), true, n.id})
		return Node{n.tx, nextId}
	}

	n.tx.SendCommand(SelectorAllSubCommand{nextId, selector, false, n.id})
	return Node{n.tx, nextId}
}

type ContentSubCommand struct {
	Id      uint64 `json:"id"`
	Content uint64 `json:"content"`
}

func (n Node) Content() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(ContentSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type ParentSubCommand struct {
	Parent uint64 `json:"parent"`
	Target uint64 `json:"target"`
}

func (n Node) Parent() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(ParentSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type FirstChildSubCommand struct {
	FirstChild uint64 `json:"firstChild"`
	Target     uint64 `json:"target"`
}

func (n Node) FirstChild() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(FirstChildSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type LastChildSubCommand struct {
	LastChild uint64 `json:"lastChild"`
	Target    uint64 `json:"target"`
}

func (n Node) LastChild() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(LastChildSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type NextSiblingSubCommand struct {
	NextSibling uint64 `json:"nextSibling"`
	Target      uint64 `json:"target"`
}

func (n Node) NextSibling() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(NextSiblingSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

type PrevSiblingSubCommand struct {
	PrevSibling uint64 `json:"prevSibling"`
	Target      uint64 `json:"target"`
}

func (n Node) PrevSibling() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(PrevSiblingSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

// Creation

type CloneSubCommand struct {
	Id    uint64 `json:"id"`
	Clone uint64 `json:"clone"`
}

func (n Node) Clone() Node {
	nextId := n.tx.getNextId()
	n.tx.SendCommand(CloneSubCommand{nextId, n.id})
	return Node{n.tx, nextId}
}

// InnerText and InnerHTML

type TextSubCommand struct {
	Target uint64 `json:"target"`
	Text   string `json:"text"`
}

func (n Node) Text(text string) Node {
	n.tx.SendCommand(TextSubCommand{n.id, text})
	return n
}

type HtmlSubCommand struct {
	Target uint64            `json:"target"`
	Html   template.Template `json:"html"`
}

func (n Node) Html(html template.Template) Node {
	n.tx.SendCommand(HtmlSubCommand{n.id, html})
	return n
}

// Attributes

type AttrSubCommand struct {
	Target uint64 `json:"target"`
	Attr   string `json:"attr"`
	Value  string `json:"value"`
}

func (n Node) Attr(name string, value string) Node {
	n.tx.SendCommand(AttrSubCommand{n.id, name, value})
	return n
}

type RemoveAttrSubCommand struct {
	Target     uint64 `json:"target"`
	RemoveAttr string `json:"removeAttr"`
}

func (n Node) RemoveAttr(name string) Node {
	n.tx.SendCommand(RemoveAttrSubCommand{n.id, name})
	return n
}

type AddToAttrSubCommand struct {
	Target    uint64 `json:"target"`
	AddToAttr string `json:"addToAttr"`
	Value     string `json:"value"`
}

func (n Node) AddToAttr(name string, value string) Node {
	n.tx.SendCommand(AddToAttrSubCommand{n.id, name, value})
	return n
}

type RemoveFromAttrSubCommand struct {
	Target         uint64 `json:"target"`
	RemoveFromAttr string `json:"removeFromAttr"`
	Value          string `json:"value"`
}

func (n Node) RemoveFromAttr(name string, value string) Node {
	n.tx.SendCommand(RemoveFromAttrSubCommand{n.id, name, value})
	return n
}

func (n Node) AddClass(class string) Node {
	return n.AddToAttr("class", class)
}

func (n Node) RemoveClass(class string) Node {
	return n.RemoveFromAttr("class", class)
}

// Node manipulation

type RemoveSubCommand struct {
	Remove uint64 `json:"remove"`
}

func (n Node) Remove() {
	n.tx.SendCommand(RemoveSubCommand{n.id})
}

type InsertNodeBeforeSubCommand struct {
	Parent           uint64 `json:"parent"`
	Ref              uint64 `json:"ref"`
	InsertNodeBefore uint64 `json:"insertNodeBefore"`
}

func (n Node) InsertNodeBefore(child Node, ref Node) Node {
	n.tx.SendCommand(InsertNodeBeforeSubCommand{n.id, ref.id, child.id})
	return n
}

type InsertBeforeSubCommand struct {
	Parent       uint64            `json:"parent"`
	Ref          uint64            `json:"ref"`
	InsertBefore template.Template `json:"insertBefore"`
}

func (n Node) InsertBefore(child template.Template, ref Node) Node {
	n.tx.SendCommand(InsertBeforeSubCommand{n.id, ref.id, child})
	return n
}

type AppendNodeSubCommand struct {
	Parent     uint64 `json:"parent"`
	AppendNode uint64 `json:"appendNode"`
}

func (n Node) AppendNode(child Node) Node {
	n.tx.SendCommand(AppendNodeSubCommand{n.id, child.id})
	return n
}

type AppendSubCommand struct {
	Parent uint64            `json:"parent"`
	Append template.Template `json:"append"`
}

func (n Node) Append(child template.Template) Node {
	n.tx.SendCommand(AppendSubCommand{n.id, child})
	return n
}

// Misc

type WaitSubCommand struct {
	Target  uint64 `json:"target"`
	Wait    string `json:"wait"`
	Timeout uint64 `json:"timeout,omitempty"`
}

func (n Node) Wait(event string, timeout uint64) Node {
	n.tx.SendCommand(WaitSubCommand{n.id, event, timeout})
	return n
}
