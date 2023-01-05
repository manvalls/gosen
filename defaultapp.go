package gosen

import "net/http"

var DefaultApp = NewApp()

func Handler(f func(p *Page, r *http.Request)) http.Handler {
	return DefaultApp.Handler(f)
}

func Handle(pattern string, f func(p *Page, r *http.Request)) {
	http.Handle(pattern, Handler(f))
}
