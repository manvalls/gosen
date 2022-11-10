package gosen

import "sync"

type transactionBuffer struct {
	commands []interface{}
	mux      sync.Mutex
}

func (b *transactionBuffer) sendCommand(id uint, command interface{}) {
	b.mux.Lock()
	defer b.mux.Unlock()

	switch c := command.(type) {
	case transactionCommand:
		b.commands = append(b.commands, c.commands...)
	default:
		b.commands = append(b.commands, command)
	}
}

type Transaction struct {
	Element
	buffer *transactionBuffer
}

func (e Element) Tx() Transaction {
	nextId := *e.nextId
	*e.nextId++

	b := &transactionBuffer{
		commands: make([]interface{}, 0),
		mux:      sync.Mutex{},
	}

	return Transaction{Element{nextId, e.nextId, e.mux, b}, b}
}

func (t Transaction) Commit() {
	t.sender.sendCommand(t.id, transactionCommand{t.buffer.commands})
}
