package gosen

import (
	"net/http"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/selectorcache"
)

type App struct {
	Hydrate       bool
	PrefetchRuns  bool
	Version       string
	GetRunHandler func(url string) http.Handler
	selectorCache *selectorcache.SelectorCache
}

func defaultGetRunHandler(url string) http.Handler {
	if url[0] == '/' {
		return http.DefaultServeMux
	}

	return nil
}

type handler struct {
	app *App
	f   func(p *Page, r *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rw, ok := w.(*commands.RunnerWriter); ok {

		h.f(&Page{
			Header:     rw.Header(),
			StatusCode: http.StatusOK,
			Routine:    rw.Routine,
		}, r)

		return
	}

	if r.URL.Query().Get("format") == "json" {
		h.serveJSON(w, r)
		return
	}

	h.serveHTML(w, r)
}

func NewApp() *App {
	return &App{
		Hydrate:       true,
		PrefetchRuns:  true,
		Version:       "",
		GetRunHandler: defaultGetRunHandler,
		selectorCache: selectorcache.New(),
	}
}

func (app *App) Handler(f func(p *Page, r *http.Request)) http.Handler {
	h := &handler{app, f}
	return h
}
