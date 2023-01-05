package gosen

import (
	"net/http"
)

type Page struct {
	Header     http.Header
	StatusCode int
	*Routine
}
