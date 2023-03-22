package gosen

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/manvalls/gosen/buffersender"
	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/htmlsender"
	"github.com/manvalls/gosen/multisender"
)

type TransactionHash struct {
	Transaction string `json:"tx"`
	Routine     uint64 `json:"routine,omitempty"`
}

type versionGetter struct {
	page *Page
}

func (v *versionGetter) Version() string {
	return v.page.Version
}

type urlPrefetcher struct {
	urlsToPrefetch map[string]bool
	mutex          *sync.Mutex
}

func (u *urlPrefetcher) PrefetchUrl(url string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.urlsToPrefetch[url] = true
}

func (h *wrappedHandler) serveHTML(w http.ResponseWriter, r *http.Request) {
	var html *htmlsender.HTMLSender
	var buffer *buffersender.BufferSender
	var sender commands.CommandSender

	html = htmlsender.NewHTMLSender(h.app.selectorCache)

	if h.app.Hydrate {
		buffer = buffersender.NewBufferSender()
		sender = multisender.NewMultiSender(buffer, html)
	} else {
		sender = html
	}

	header := w.Header()
	wg := &sync.WaitGroup{}

	p := &Page{
		Version: h.app.Version,
		Header:  header,
		writter: w,
	}

	runner := &commands.Runner{
		VersionGetter:    &versionGetter{p},
		RunHandlerGetter: h.app.RunHandlerGetter,
		BaseRequest:      r,
		Header:           header,
	}

	var urlsToPrefetch map[string]bool

	if h.app.PrefetchRuns {
		urlsToPrefetch = map[string]bool{}
		runner.UrlPrefetcher = &urlPrefetcher{urlsToPrefetch, &sync.Mutex{}}
	}

	p.Routine = commands.NewRoutine(sender, wg, runner)
	h.handler.ServeGosen(p, r)
	wg.Wait()

	if h.app.Hydrate {
		cmdList := []any{}
		for _, cmd := range buffer.Commands() {
			switch c := cmd.(type) {
			case commands.TransactionCommand:
				cmdList = append(cmdList, TransactionHash{c.Hash, c.Routine})
			default:
				cmdList = append(cmdList, cmd)
			}
		}

		hydrationData, err := json.Marshal(cmdList)
		if err == nil {
			script := "window.__GOSEN_HYDRATION__=" + string(hydrationData) + ";"

			if p.Version != "" {
				version, _ := json.Marshal(p.Version)
				script += "window.__GOSEN_PAGE_VERSION__=" + string(version) + ";"
			}

			html.PrependScript(script)
		}
	}

	for url := range urlsToPrefetch {
		html.Prefecth(url)
	}

	if h.runCache != nil {
		needsUpdating := false

		h.runCacheMux.RLock()

		for url := range urlsToPrefetch {
			if !h.runCache[url] {
				needsUpdating = true
				break
			}
		}

		h.runCacheMux.RUnlock()

		if needsUpdating {
			h.runCacheMux.Lock()

			for url := range urlsToPrefetch {
				h.runCache[url] = true
			}

			h.runCacheMux.Unlock()
		}
	}

	html.Render(w)
}
