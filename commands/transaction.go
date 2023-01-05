package commands

import (
	"sync"

	"github.com/manvalls/gosen/template"
)

type Transaction struct {
	sender   CommandSender
	commands []interface{}
	mux      *sync.Mutex
	nextId   uint
	routine  uint
}

type TransactionCommand struct {
	Transaction []any `json:"transaction"`
	Routine     uint  `json:"routine,omitempty"`
}

func (t *Transaction) getNextId() uint {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.nextId++
	return t.nextId
}

func (t *Transaction) SendCommand(command any) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.commands = append(t.commands, command)
}

func (t *Transaction) Commit() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.commands) == 0 {
		return
	}

	commands := t.commands
	t.commands = nil
	t.sender.SendCommand(TransactionCommand{commands, t.routine})
}

// Selectors

type SelectorSubCommand struct {
	Id       uint          `json:"id"`
	Selector string        `json:"selector"`
	Args     []interface{} `json:"args,omitempty"`
	Parent   uint          `json:"parent,omitempty"`
}

func (t *Transaction) S(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.SendCommand(SelectorSubCommand{nextId, selector, args, 0})
	return Node{t, nextId}
}

type SelectorAllSubCommand struct {
	Id          uint          `json:"id"`
	SelectorAll string        `json:"selectorAll"`
	Args        []interface{} `json:"args,omitempty"`
	Parent      uint          `json:"parent,omitempty"`
}

func (t *Transaction) All(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.SendCommand(SelectorAllSubCommand{nextId, selector, args, 0})
	return Node{t, nextId}
}

// Creation

type FragmentSubCommand struct {
	Id       uint              `json:"id"`
	Fragment template.Template `json:"fragment"`
}

func (t *Transaction) Fragment(fragment template.Template) Node {
	nextId := t.getNextId()
	t.SendCommand(FragmentSubCommand{nextId, fragment})
	return Node{t, nextId}
}
