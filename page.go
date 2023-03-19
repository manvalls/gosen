package gosen

import (
	"net/http"
)

type Page struct {
	Version    string
	Header     http.Header
	StatusCode int
	*Routine
}
