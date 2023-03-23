package gosen

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Page struct {
	Version string
	Header  http.Header
	*Routine
	writter http.ResponseWriter
	sseMux  *sync.Mutex
}

func (p *Page) WriteHeader(statusCode int) {
	p.writter.WriteHeader(statusCode)
}

type Event struct {
	Event string
	Data  string
	Id    string
	Retry int
}

func (p *Page) SendEvent(e Event) {
	if p.sseMux == nil {
		return
	}

	p.sseMux.Lock()
	defer p.sseMux.Unlock()

	if e.Event != "" {
		p.writter.Write([]byte("event: " + strings.ReplaceAll(e.Event, "\n", "") + "\n"))
	}

	if e.Data != "" {
		p.writter.Write([]byte("data: " + strings.ReplaceAll(e.Data, "\n", "") + "\n"))
	}

	if e.Id != "" {
		p.writter.Write([]byte("id: " + strings.ReplaceAll(e.Id, "\n", "") + "\n"))
	}

	if e.Retry != 0 {
		p.writter.Write([]byte("retry: " + strconv.Itoa(e.Retry) + "\n"))
	}

	p.writter.Write([]byte("\n"))

	if flusher, ok := p.writter.(http.Flusher); ok {
		flusher.Flush()
	}
}
