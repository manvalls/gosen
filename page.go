package gosen

import (
	"net/http"
)

type Page struct {
	Version string
	Header  http.Header
	*Routine
	writter http.ResponseWriter
}

func (p *Page) WriteHeader(statusCode int) {
	p.writter.WriteHeader(statusCode)
}
