package gosen

import (
	"encoding/json"
	"net/http"

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

	p := &Page{
		Header:     w.Header(),
		StatusCode: http.StatusOK,
		Routine:    commands.NewRoutine(sender),
	}

	h.f(p, r)

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
			tx := p.Tx()
			head := tx.S("head")

			head.InsertBefore(
				Raw(`<script type="text/javascript">window.__GOSEN_HYDRATION__ = `+string(hydrationData)+`;</script>`),
				head.FirstChild(),
			)

			tx.Commit()
		}
	}

	w.WriteHeader(p.StatusCode)
	html.Render(w)
}
