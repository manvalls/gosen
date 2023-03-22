package gosen

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/manvalls/gosen/buffersender"
	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/util"
)

type JSONResponse struct {
	Commands []any  `json:"commands"`
	Version  string `json:"version,omitempty"`
}

func (h *wrappedHandler) serveJSON(w http.ResponseWriter, r *http.Request) {
	if h.runCache != nil {
		h.runCacheMux.RLock()

		for url := range h.runCache {
			w.Header().Add("Link", "<"+url+">; rel=preload; as=fetch")
		}

		h.runCacheMux.RUnlock()

		w.WriteHeader(http.StatusEarlyHints)
	}

	buffer := buffersender.NewBufferSender()
	header := w.Header()
	wg := &sync.WaitGroup{}

	p := &Page{
		Version: h.app.Version,
		Header:  header,
		Routine: commands.NewRoutine(buffer, wg, nil),
		writter: w,
	}

	h.handler.ServeGosen(p, r)
	wg.Wait()

	cmd := buffer.GetCommands()

	if h.runCache != nil {
		urls := []string{}

		query := "format=json"
		if p.Version != "" {
			query += "&version=" + p.Version
		}

		for _, c := range cmd {
			if c, ok := c.(commands.RunCommand); ok && !c.Dynamic {
				urls = append(urls, util.AddToQuery(c.Run, query))
			}
		}

		somethingToChange := false

		h.runCacheMux.RLock()

		for _, url := range urls {
			if !h.runCache[url] {
				somethingToChange = true
				break
			}
		}

		h.runCacheMux.RUnlock()

		if somethingToChange {
			h.runCacheMux.Lock()

			for _, url := range urls {
				h.runCache[url] = true
			}

			h.runCacheMux.Unlock()
		}
	}

	data, _ := json.Marshal(JSONResponse{
		Commands: cmd,
		Version:  p.Version,
	})

	w.Write(data)
}
