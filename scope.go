package gosen

import (
	"sync"
)

type Scope interface {
	Run(url string) Scope
	RunOnce(url string) Scope
	Scope() Scope
	Tx() *Transaction
	Commit()
}

type scope struct {
	sender   commandSender
	commands []interface{}
	mux      sync.Mutex
}

func (s *scope) sendCommand(command interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.commands = append(s.commands, command)
}

func (s *scope) Run(url string) Scope {
	s.sendCommand(runCommand{url})
	return s
}

func (s *scope) RunOnce(url string) Scope {
	s.sendCommand(onceCommand{url})
	return s
}

func (s *scope) Scope() Scope {
	return &scope{s, nil, sync.Mutex{}}
}

func (s *scope) Tx() *Transaction {
	return &Transaction{s, nil, sync.Mutex{}, 0}
}

func (s *scope) Commit() {
	s.mux.Lock()
	defer s.mux.Unlock()

	if len(s.commands) == 0 {
		return
	}

	commands := s.commands
	s.commands = nil
	s.sender.sendCommand(scopeCommand{commands})
}
