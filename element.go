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

// TODO: attributes, append, remove, insertBefore, replace, etc.

// Scripts

func (e Element) LoadScript(scriptURL string) Element {
	e.sender.sendCommand(e.id, loadScriptCommand{scriptURL, false})
	return e
}

func (e Element) LoadScriptAsync(scriptURL string) Element {
	e.sender.sendCommand(e.id, loadScriptCommand{scriptURL, true})
	return e
}
