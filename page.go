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

func (p *Page) Run(url string) *Page {
	p.sender.sendCommand(runCommand{url})
	return p
}

func (p *Page) RunAsync(url string) *Page {
	p.sender.sendCommand(asyncCommand{url})
	return p
}

func (p *Page) RunOnce(url string) *Page {
	p.sender.sendCommand(onceCommand{url})
	return p
}

func (p *Page) RunOnceAsync(url string) *Page {
	p.sender.sendCommand(onceAsyncCommand{url})
	return p
}

func (p *Page) Tx() *Transaction {
	return &Transaction{p, nil, sync.Mutex{}, 0}
}
