package gosen

import (
	"net/http"
	"sync"
	"time"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/ssesender"
)

func (h *wrappedHandler) serveSSE(w http.ResponseWriter, r *http.Request) {
	h.sendEarlyHints(w)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if h.app.SSEKeepAlive > 0 {
		go func() {
			for {
				select {
				case <-r.Context().Done():
					return
				case <-time.After(time.Duration(h.app.SSEKeepAlive) * time.Second):
					w.Write([]byte(":keepalive\n\n"))
					flusher.Flush()
				}
			}
		}()
	}

	p := &Page{
		Version: h.app.Version,
		Header:  w.Header(),
		writter: w,
	}

	sender := &ssesender.SSESender{
		Writter:       w,
		Flusher:       flusher,
		VersionGetter: &versionGetter{p},
		RunList:       []string{},
	}

	wg := &sync.WaitGroup{}
	p.Routine = commands.NewRoutine(sender, wg, nil)

	h.handler.ServeGosen(p, r)
	wg.Wait()

	h.cacheRuns(p.Version, sender.RunList)
}
