package gosen

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/manvalls/gosen/selectorcache"
)

type Config struct {
	SSEKeepAlive int
	noHydrate    bool
	Version      string
	VersionFunc  func() string
	MapRunURL    func(string) string
}

type gosenHandler struct {
	config        Config
	handler       http.Handler
	selectorCache *selectorcache.SelectorCache
}

type gosenContext struct {
	version string
	mux     sync.Mutex
	routine *Routine
	*gosenHandler
	http.ResponseWriter

	// WaitGroup for open pages
	pending sync.WaitGroup

	// Channel that gets closed when the request finishes
	done chan struct{}
}

var gosenContextKey = &struct{}{}

func (h *gosenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	version := h.config.Version
	if h.config.VersionFunc != nil {
		version = h.config.VersionFunc()
	}

	clientVersion := r.URL.Query().Get("version")
	if version != "" && clientVersion != "" && version != clientVersion {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"VERSION_MISMATCH","serverVersion":"` + strings.ReplaceAll(version, `"`, "") + `","clientVersion":"` + strings.ReplaceAll(clientVersion, `"`, "") + `"}`))
		return
	}

	c := &gosenContext{
		version:        version,
		gosenHandler:   h,
		ResponseWriter: w,
		mux:            sync.Mutex{},
		pending:        sync.WaitGroup{},
		done:           make(chan struct{}),
	}

	gosenRequest := r.Clone(context.WithValue(r.Context(), gosenContextKey, c))
	h.handler.ServeHTTP(w, gosenRequest)
	close(c.done)
	c.pending.Wait()
}

var App = AppWithConfig(Config{})

func AppWithConfig(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &gosenHandler{
			config:        config,
			handler:       next,
			selectorCache: selectorcache.New(),
		}
	}
}
