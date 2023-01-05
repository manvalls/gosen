package gosen

import (
	"net/http"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/htmlsender"
)

func (h *handler) serveHTML(w http.ResponseWriter, r *http.Request) {
	sender := htmlsender.NewHTMLSender(h.app.selectorCache)
	p := &Page{
		Header:     w.Header(),
		StatusCode: http.StatusOK,
		Routine:    commands.NewRoutine(sender),
	}

	h.f(p, r)

	// TODO: add hashed list of commands

	w.WriteHeader(p.StatusCode)
	sender.Render(w)
}
