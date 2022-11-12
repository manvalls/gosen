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
	e.sender.sendCommand(selectorCommand{e.id, nextId, selector, args})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) All(selector string, args ...interface{}) Node {
	nextId := e.getNextId()
	e.sender.sendCommand(selectorAllCommand{e.id, nextId, selector, args})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) Content() Node {
	nextId := e.getNextId()
	e.sender.sendCommand(contentCommand{e.id, nextId})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

// Creation

func (e Node) Fragment(template Template) Node {
	nextId := e.getNextId()
	e.sender.sendCommand(fragmentCommand{nextId, template})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

func (e Node) Clone() Node {
	nextId := e.getNextId()
	e.sender.sendCommand(cloneCommand{e.id, nextId})
	return Node{nextId, e.nextId, e.mux, e.sender}
}

// InnerText and InnerHTML

func (e Node) Text(text string) Node {
	e.sender.sendCommand(textCommand{e.id, text})
	return e
}

func (e Node) HTML(html Template) Node {
	e.sender.sendCommand(htmlCommand{e.id, html})
	return e
}

// Attributes

func (e Node) Attr(name string, value string) Node {
	e.sender.sendCommand(attrCommand{e.id, name, value})
	return e
}

func (e Node) RmAttr(name string) Node {
	e.sender.sendCommand(rmAttrCommand{e.id, name})
	return e
}

func (e Node) AddToAttr(name string, value string) Node {
	e.sender.sendCommand(addToAttrCommand{e.id, name, value})
	return e
}

func (e Node) RmFromAttr(name string, value string) Node {
	e.sender.sendCommand(rmFromAttrCommand{e.id, name, value})
	return e
}

func (e Node) AddClass(name string) Node {
	e.sender.sendCommand(addClassCommand{e.id, name})
	return e
}

func (e Node) RmClass(name string) Node {
	e.sender.sendCommand(rmClassCommand{e.id, name})
	return e
}

// Node manipulation

func (e Node) Remove() Node {
	e.sender.sendCommand(removeCommand{e.id})
	return e
}

func (e Node) Empty() Node {
	e.sender.sendCommand(emptyCommand{e.id})
	return e
}

func (e Node) ReplaceWith(otherNode Node) Node {
	e.sender.sendCommand(replaceWithCommand{e.id, otherNode.id})
	return e
}

func (e Node) InsertBefore(child Node, ref Node) Node {
	e.sender.sendCommand(insertBeforeCommand{e.id, child.id, ref.id})
	return e
}

func (e Node) InsertAfter(child Node, ref Node) Node {
	e.sender.sendCommand(insertAfterCommand{e.id, child.id, ref.id})
	return e
}

func (e Node) Append(child Node) Node {
	e.sender.sendCommand(appendCommand{e.id, child.id})
	return e
}

func (e Node) Prepend(child Node) Node {
	e.sender.sendCommand(prependCommand{e.id, child.id})
	return e
}

// Misc

func (e Node) Wait(event string, timeout uint) Node {
	e.sender.sendCommand(waitCommand{e.id, event, timeout})
	return e
}

func (e Node) Run(url string) Node {
	e.sender.sendCommand(runCommand{url})
	return e
}

func (e Node) Listen(url string) Node {
	e.sender.sendCommand(listenCommand{url})
	return e
}

func (e Node) Async(url string) Node {
	e.sender.sendCommand(asyncCommand{url})
	return e
}

func (e Node) Defer(url string) Node {
	e.sender.sendCommand(deferCommand{url})
	return e
}

func (e Node) Once(url string) Node {
	e.sender.sendCommand(onceCommand{url})
	return e
}
