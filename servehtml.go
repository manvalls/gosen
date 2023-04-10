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

func (h *wrappedHandler) serveHTML(w http.ResponseWriter, r *http.Request, version string) {
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
		Header:  header,
		writter: w,
	}

	runner := &commands.Runner{
		BaseRequest: r,
		Header:      header,
	}

	p.Routine = commands.NewRoutine(sender, wg, runner)
	h.handler.ServeGosen(p, r)
	wg.Wait()

	if h.app.Hydrate {
		cmdList := []any{}
		for _, cmd := range buffer.Commands() {
			switch c := cmd.(type) {
			case commands.TransactionCommand:
				cmdList = append(cmdList, commands.TransactionCommand{
					Hash:    c.Hash,
					Routine: c.Routine,
					Once:    c.Once,
				})
			default:
				cmdList = append(cmdList, cmd)
			}
		}

		hydrationData, err := json.Marshal(cmdList)
		if err == nil {
			script := "window.__GOSEN_STATE__=" + string(hydrationData) + ";"
			if version != "" {
				v, _ := json.Marshal(version)
				script += "window.__GOSEN_PAGE_VERSION__=" + string(v) + ";"
			}

			html.PrependScript(script)
		}
	}

	html.Render(w)
}
