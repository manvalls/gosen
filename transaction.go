package gosen

import "sync"

type transactionBuffer struct {
	commands []interface{}
	mux      sync.Mutex
}

func (b *transactionBuffer) sendCommand(command interface{}) {
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
	Node
	buffer *transactionBuffer
}

func (e Node) Tx() Transaction {
	b := &transactionBuffer{
		commands: make([]interface{}, 0),
		mux:      sync.Mutex{},
	}

	return Transaction{Node{e.id, e.nextId, e.mux, b}, b}
}

func (t Transaction) Commit() {
	t.sender.sendCommand(transactionCommand{t.buffer.commands})
}
