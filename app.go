package gosen

import (
	"net/http"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/selectorcache"
)

type DefaultRunHandlerGetter struct{}

func (d *DefaultRunHandlerGetter) RunHandler(url string) http.Handler {
	if url[0] == '/' {
		return http.DefaultServeMux
	}

	return nil
}

type App struct {
	Hydrate       bool
	PrefetchRuns  bool
	Version       string
	selectorCache *selectorcache.SelectorCache
	commands.RunHandlerGetter
}

type Handler interface {
	ServeGosen(p *Page, r *http.Request)
}

type wrappedHandler struct {
	app     *App
	handler Handler
}

func (h *wrappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rw, ok := w.(*commands.RunnerWriter); ok {

		h.handler.ServeGosen(&Page{
			Version: h.app.Version,
			Header:  rw.Header(),
			Routine: rw.Routine,
			writter: w,
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
		Hydrate:          true,
		PrefetchRuns:     true,
		Version:          "",
		RunHandlerGetter: &DefaultRunHandlerGetter{},
		selectorCache:    selectorcache.New(),
	}
}

func (app *App) Wrap(h Handler) http.Handler {
	return &wrappedHandler{app, h}
}

type funcHandler struct {
	f func(p *Page, r *http.Request)
}

func (h *funcHandler) ServeGosen(p *Page, r *http.Request) {
	h.f(p, r)
}

func (app *App) WrapFunc(f func(p *Page, r *http.Request)) http.Handler {
	return app.Wrap(&funcHandler{f})
}
