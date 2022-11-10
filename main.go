package gosen

import "net/http"

type App struct {
}

type Page struct {
	Node
	Header     http.Header
	StatusCode int
}

type handler struct {
}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Handler(func(p Page, r *http.Request)) http.Handler {
	// TODO
	h := &handler{}
	return h
}
