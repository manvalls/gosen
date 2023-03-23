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

	flusher, _ := w.(http.Flusher)

	if flusher != nil && h.app.SSEKeepAlive > 0 {
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

	mux := &sync.Mutex{}

	sender := &ssesender.SSESender{
		Mux:           mux,
		Writter:       w,
		Flusher:       flusher,
		VersionGetter: &versionGetter{p},
		RunList:       []string{},
	}

	wg := &sync.WaitGroup{}
	p.Routine = commands.NewRoutine(sender, wg, nil)
	p.sseMux = mux

	h.handler.ServeGosen(p, r)
	wg.Wait()

	h.cacheRuns(p.Version, sender.RunList)
}
