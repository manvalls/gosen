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

func (h *handler) serveHTML(w http.ResponseWriter, r *http.Request) {
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

	var p *Page

	p = &Page{
		Version:    h.app.Version,
		Header:     header,
		StatusCode: http.StatusOK,
		Routine: commands.NewRoutine(sender, wg, &commands.Runner{
			Version: func() string {
				return p.Version
			},
			GetRunHandler: h.app.GetRunHandler,
			BaseRequest:   r,
			Header:        header,
		}),
	}

	h.f(p, r)
	wg.Wait()

	if h.app.Hydrate {
		cmdList := []any{}
		for _, cmd := range buffer.GetCommands() {
			switch c := cmd.(type) {
			case commands.TransactionCommand:
				cmdList = append(cmdList, TransactionHash{c.Hash, c.Routine})
			default:
				cmdList = append(cmdList, cmd)
			}
		}

		hydrationData, err := json.Marshal(cmdList)
		if err == nil {
			version, _ := json.Marshal(p.Version)
			html.PrependScript("window.__GOSEN_HYDRATION__=" + string(hydrationData) + ";window.__GOSEN_PAGE_VERSION__=" + string(version) + ";")
		}
	}

	w.WriteHeader(p.StatusCode)
	html.Render(w)
}
