package gosen

import "sync"

type Transaction struct {
	page     *Page
	commands []interface{}
	mux      sync.Mutex
	nextId   uint
}

func (t *Transaction) getNextId() uint {
	t.mux.Lock()
	defer t.mux.Unlock()
	nextId := t.nextId
	t.nextId++
	return nextId
}

func (t *Transaction) addCommand(command interface{}) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.commands = append(t.commands, command)
}

func (t *Transaction) Commit() {
	t.mux.Lock()
	defer t.mux.Unlock()

	commands := t.commands
	t.commands = nil
	t.page.sender.sendCommand(transactionCommand{commands})
}

// Selectors

func (t *Transaction) S(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.addCommand(rootSelectorCommand{nextId, selector, args})
	return Node{t, nextId}
}

func (t *Transaction) All(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.addCommand(rootSelectorAllCommand{nextId, selector, args})
	return Node{t, nextId}
}

// Creation

func (t *Transaction) Fragment(template Template) Node {
	nextId := t.getNextId()
	t.addCommand(fragmentCommand{nextId, template})
	return Node{t, nextId}
}
