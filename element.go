package gosen

type Element struct {
	id     uint
	nextId *uint
	sender commandSender
}

type Transaction struct {
	Element
}

// Transactions

func (e Element) Tx() Transaction {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, transactionCommand{nextId})
	return Transaction{Element{nextId, e.nextId, e.sender}}
}

func (t Transaction) Commit() {
	t.sender.sendCommand(t.id, commitCommand{})
}

// Selectors

func (e Element) S(selector string, args ...interface{}) Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, selectorCommand{nextId, selector, args})
	return Element{nextId, e.nextId, e.sender}
}

func (e Element) All(selector string, args ...interface{}) Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, selectorAllCommand{nextId, selector, args})
	return Element{nextId, e.nextId, e.sender}
}

func (e Element) Content() Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, contentCommand{nextId})
	return Element{nextId, e.nextId, e.sender}
}

// Creation

func (e Element) Create(template Template) Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, createCommand{nextId, template})
	return Element{nextId, e.nextId, e.sender}
}

func (e Element) Fragment(template Template) Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, fragmentCommand{nextId, template})
	return Element{nextId, e.nextId, e.sender}
}

func (e Element) Clone() Element {
	nextId := *e.nextId
	*e.nextId++
	e.sender.sendCommand(e.id, cloneCommand{nextId})
	return Element{nextId, e.nextId, e.sender}
}

// InnerText and InnerHTML

func (e Element) Text(text string) Element {
	e.sender.sendCommand(e.id, textCommand{text})
	return e
}

func (e Element) HTML(html Template) Element {
	e.sender.sendCommand(e.id, htmlCommand{html})
	return e
}

// Attributes

func (e Element) Attr(name string, value string) Element {
	e.sender.sendCommand(e.id, attrCommand{name, value})
	return e
}

func (e Element) RmAttr(name string) Element {
	e.sender.sendCommand(e.id, rmAttrCommand{name})
	return e
}

func (e Element) AddToAttr(name string, value string) Element {
	e.sender.sendCommand(e.id, addToAttrCommand{name, value})
	return e
}

func (e Element) RmFromAttr(name string, value string) Element {
	e.sender.sendCommand(e.id, rmFromAttrCommand{name, value})
	return e
}

func (e Element) AddClass(name string) Element {
	e.sender.sendCommand(e.id, addClassCommand{name})
	return e
}

func (e Element) RmClass(name string) Element {
	e.sender.sendCommand(e.id, rmClassCommand{name})
	return e
}

// Element manipulation

func (e Element) Remove() Element {
	e.sender.sendCommand(e.id, removeCommand{})
	return e
}

func (e Element) Empty() Element {
	e.sender.sendCommand(e.id, emptyCommand{})
	return e
}

func (e Element) ReplaceWith(otherElement Element) Element {
	e.sender.sendCommand(e.id, replaceWithCommand{otherElement.id})
	return e
}

func (e Element) InsertBefore(child Element, ref Element) Element {
	e.sender.sendCommand(e.id, insertBeforeCommand{child.id, ref.id})
	return e
}

func (e Element) InsertAfter(child Element, ref Element) Element {
	e.sender.sendCommand(e.id, insertAfterCommand{child.id, ref.id})
	return e
}

func (e Element) Append(child Element) Element {
	e.sender.sendCommand(e.id, appendCommand{child.id})
	return e
}

func (e Element) Prepend(child Element) Element {
	e.sender.sendCommand(e.id, prependCommand{child.id})
	return e
}

// Misc

func (e Element) Wait(event string) Element {
	e.sender.sendCommand(e.id, waitCommand{event})
	return e
}

func (e Element) Run(url string) Element {
	e.sender.sendCommand(e.id, runCommand{url})
	return e
}

func (e Element) Listen(url string) Element {
	e.sender.sendCommand(e.id, listenCommand{url})
	return e
}

func (e Element) Async(url string) Element {
	e.sender.sendCommand(e.id, asyncCommand{url})
	return e
}

func (e Element) Defer(url string) Element {
	e.sender.sendCommand(e.id, deferCommand{url})
	return e
}

func (e Element) Once(url string) Element {
	e.sender.sendCommand(e.id, onceCommand{url})
	return e
}
