package gosen

import "net/http"

var DefaultApp = NewApp()

func WrapFunc(f func(p *Page, r *http.Request)) http.Handler {
	return DefaultApp.WrapFunc(f)
}

func Wrap(h Handler) http.Handler {
	return DefaultApp.Wrap(h)
}

func HandleFunc(pattern string, f func(p *Page, r *http.Request)) {
	http.Handle(pattern, WrapFunc(f))
}

func Handle(pattern string, h Handler) {
	http.Handle(pattern, Wrap(h))
}
