package gosen

import "sync"

type Node struct {
	id     uint
	nextId *uint
	mux    *sync.Mutex
	sender commandSender
}

func (e Node) getNextId() uint {
	e.mux.Lock()
	defer e.mux.Unlock()
	nextId := *e.nextId
	*e.nextId++
	return nextId
}

// Selectors

func (e Node) S(selector string, args ...interface{}) Node {
	nextId := e.getNextId()
	e.sender.sendCommand(e.id, selectorCommand{nextId, selector, args})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) All(selector string, args ...interface{}) Node {
	nextId := e.getNextId()
	e.sender.sendCommand(e.id, selectorAllCommand{nextId, selector, args})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) Content() Node {
	nextId := e.getNextId()
	e.sender.sendCommand(e.id, contentCommand{nextId})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

// Creation

func (e Node) Fragment(template Template) Node {
	nextId := e.getNextId()
	e.sender.sendCommand(e.id, fragmentCommand{nextId, template})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) Clone() Node {
	nextId := e.getNextId()
	e.sender.sendCommand(e.id, cloneCommand{nextId})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

// InnerText and InnerHTML

func (e Node) Text(text string) Node {
	e.sender.sendCommand(e.id, textCommand{text})
	return e
}

func (e Node) HTML(html Template) Node {
	e.sender.sendCommand(e.id, htmlCommand{html})
	return e
}

// Attributes

func (e Node) Attr(name string, value string) Node {
	e.sender.sendCommand(e.id, attrCommand{name, value})
	return e
}

func (e Node) RmAttr(name string) Node {
	e.sender.sendCommand(e.id, rmAttrCommand{name})
	return e
}

func (e Node) AddToAttr(name string, value string) Node {
	e.sender.sendCommand(e.id, addToAttrCommand{name, value})
	return e
}

func (e Node) RmFromAttr(name string, value string) Node {
	e.sender.sendCommand(e.id, rmFromAttrCommand{name, value})
	return e
}

func (e Node) AddClass(name string) Node {
	e.sender.sendCommand(e.id, addClassCommand{name})
	return e
}

func (e Node) RmClass(name string) Node {
	e.sender.sendCommand(e.id, rmClassCommand{name})
	return e
}

// Node manipulation

func (e Node) Remove() Node {
	e.sender.sendCommand(e.id, removeCommand{})
	return e
}

func (e Node) Empty() Node {
	e.sender.sendCommand(e.id, emptyCommand{})
	return e
}

func (e Node) ReplaceWith(otherNode Node) Node {
	e.sender.sendCommand(e.id, replaceWithCommand{otherNode.id})
	return e
}

func (e Node) InsertBefore(child Node, ref Node) Node {
	e.sender.sendCommand(e.id, insertBeforeCommand{child.id, ref.id})
	return e
}

func (e Node) InsertAfter(child Node, ref Node) Node {
	e.sender.sendCommand(e.id, insertAfterCommand{child.id, ref.id})
	return e
}

func (e Node) Append(child Node) Node {
	e.sender.sendCommand(e.id, appendCommand{child.id})
	return e
}

func (e Node) Prepend(child Node) Node {
	e.sender.sendCommand(e.id, prependCommand{child.id})
	return e
}

// Misc

func (e Node) Wait(event string) Node {
	e.sender.sendCommand(e.id, waitCommand{event})
	return e
}

func (e Node) Run(url string) Node {
	e.sender.sendCommand(e.id, runCommand{url})
	return e
}

func (e Node) Listen(url string) Node {
	e.sender.sendCommand(e.id, listenCommand{url})
	return e
}

func (e Node) Async(url string) Node {
	e.sender.sendCommand(e.id, asyncCommand{url})
	return e
}

func (e Node) Defer(url string) Node {
	e.sender.sendCommand(e.id, deferCommand{url})
	return e
}

func (e Node) Once(url string) Node {
	e.sender.sendCommand(e.id, onceCommand{url})
	return e
}
