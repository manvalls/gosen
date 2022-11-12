package gosen

import "sync"

type cleanupBuffer struct {
	commands []interface{}
	mux      sync.Mutex
}

func (b *cleanupBuffer) sendCommand(command interface{}) {
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
	b := &cleanupBuffer{
		commands: make([]interface{}, 0),
		mux:      sync.Mutex{},
	}

	return Cleanup{Node{e.id, e.nextId, e.mux, b}, b}
}

func (t Cleanup) Commit() {
	t.sender.sendCommand(cleanupCommand{t.buffer.commands})
}
