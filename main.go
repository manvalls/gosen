package gosen

import (
	"net/http"

	"github.com/manvalls/gosen/selectorcache"
)

type App struct {
	selectorCache *selectorcache.SelectorCache
}

type handler struct {
	app *App
	f   func(p *Page, r *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.serveHTML(w, r)
}

func NewApp() *App {
	return &App{
		selectorCache: selectorcache.New(),
	}
}

func (app *App) Handler(f func(p *Page, r *http.Request)) http.Handler {
	h := &handler{app, f}
	return h
}
