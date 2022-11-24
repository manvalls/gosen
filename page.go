package gosen

import (
	"net/http"
	"sync"
)

type Page struct {
	Header     http.Header
	StatusCode int

	sender commandSender
}

func (p *Page) Run(url string) Scope {
	p.sender.sendCommand(runCommand{url})
	return p
}

func (p *Page) RunOnce(url string) Scope {
	p.sender.sendCommand(onceCommand{url})
	return p
}

func (p *Page) Tx() *Transaction {
	return &Transaction{p.sender, nil, sync.Mutex{}, 0}
}

func (p *Page) Scope() Scope {
	return &scope{p.sender, nil, sync.Mutex{}}
}

func (p *Page) Commit() {
	// Do nothing
}
