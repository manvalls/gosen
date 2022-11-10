package gosen

import "sync"

type cleanupBuffer struct {
	commands []interface{}
	mux      sync.Mutex
}

func (b *cleanupBuffer) sendCommand(id uint, command interface{}) {
	b.mux.Lock()
	defer b.mux.Unlock()

	switch c := command.(type) {
	case transactionCommand:
		b.commands = append(b.commands, c.commands...)
	case cleanupCommand:
		b.commands = append(b.commands, c.commands...)
	default:
		b.commands = append(b.commands, command)
	}
}

type Cleanup struct {
	Node
	buffer *cleanupBuffer
}

func (e Node) Cleanup() Cleanup {
	nextId := *e.nextId
	*e.nextId++

	b := &cleanupBuffer{
		commands: make([]interface{}, 0),
		mux:      sync.Mutex{},
	}

	return Cleanup{Node{nextId, e.nextId, e.mux, b}, b}
}

func (t Cleanup) Commit() {
	t.sender.sendCommand(t.id, cleanupCommand{t.buffer.commands})
}
