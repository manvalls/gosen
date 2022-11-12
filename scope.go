package gosen

import "sync"

type scopeBuffer struct {
	commands []interface{}
	mux      sync.Mutex
}

func (b *scopeBuffer) sendCommand(command interface{}) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.commands = append(b.commands, command)
}

type Scope struct {
	Node
	id     string
	buffer *scopeBuffer
}

func (e Node) Scope(id string) Scope {
	b := &scopeBuffer{
		commands: make([]interface{}, 0),
		mux:      sync.Mutex{},
	}

	return Scope{Node{e.id, e.nextId, e.mux, b}, id, b}
}

func (t Scope) Commit() {
	t.sender.sendCommand(scopeCommand{t.buffer.commands, t.id})
}
