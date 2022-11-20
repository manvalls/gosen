package gosen

import (
	"sync"
)

type Async struct {
	sender   commandSender
	commands []interface{}
	mux      sync.Mutex
}

func (a *Async) addCommand(command interface{}) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.commands = append(a.commands, command)
}

func (a *Async) Commit() {
	a.mux.Lock()
	defer a.mux.Unlock()

	commands := a.commands
	a.commands = nil
	a.sender.sendCommand(asyncCommand{commands})
}

func (a *Async) Run(url string) *Async {
	a.addCommand(runCommand{url})
	return a
}

func (a *Async) RunOnce(url string) *Async {
	a.addCommand(onceCommand{url})
	return a
}

func (a *Async) Tx() *Transaction {
	return &Transaction{a.sender, nil, sync.Mutex{}, 0}
}
