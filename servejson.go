package gosen

import (
	"net/http"
	"sync"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/jsonsender"
)

func (h *wrappedHandler) serveJSON(w http.ResponseWriter, r *http.Request) {
	h.sendEarlyHints(w)

	p := &Page{
		Version: h.app.Version,
		Header:  w.Header(),
		writter: w,
	}

	sender := &jsonsender.JSONSender{
		Writter:       w,
		VersionGetter: &versionGetter{p},
		RunList:       []string{},
	}

	wg := &sync.WaitGroup{}
	p.Routine = commands.NewRoutine(sender, wg, nil)

	h.handler.ServeGosen(p, r)
	wg.Wait()
	sender.End()

	h.cacheRuns(p.Version, sender.RunList)
}
