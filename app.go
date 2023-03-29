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
	SSEKeepAlive  int
	Hydrate       bool
	PrefetchRuns  bool
	Version       string
	selectorCache *selectorcache.SelectorCache
	commands.RunHandlerGetter
	commands.VersionGetter
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
			Header:  rw.Header(),
			Routine: rw.Routine,
			writter: w,
		}, r)

		return
	}

	if r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		h.serveJSON(w, r)
		return
	}

	if r.URL.Query().Get("format") == "sse" {
		w.Header().Set("Content-Type", "text/event-stream")
		h.serveSSE(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.serveHTML(w, r)
}

type versionGetter struct {
	app *App
}

func (v *versionGetter) Version() string {
	return v.app.Version
}

func NewApp() *App {
	app := &App{
		SSEKeepAlive:     15,
		Hydrate:          true,
		PrefetchRuns:     true,
		Version:          "",
		RunHandlerGetter: &DefaultRunHandlerGetter{},
		selectorCache:    selectorcache.New(),
	}

	app.VersionGetter = &versionGetter{app}
	return app
}

func (app *App) Wrap(h Handler) http.Handler {
	wh := &wrappedHandler{
		app:     app,
		handler: h,
	}
	return wh
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
