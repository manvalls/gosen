package gosen

import (
	"encoding/json"
	"net/http"

	"github.com/manvalls/gosen/buffersender"
	"github.com/manvalls/gosen/commands"
)

func (h *handler) serveJSON(w http.ResponseWriter, r *http.Request) {
	buffer := buffersender.NewBufferSender()
	header := w.Header()

	p := &Page{
		Header:     header,
		StatusCode: http.StatusOK,
		Routine:    commands.NewRoutine(buffer),
	}

	h.f(p, r)

	header.Add("vary", "gosen-accept")
	w.WriteHeader(p.StatusCode)
	data, _ := json.Marshal(buffer.GetCommands())
	w.Write(data)
}
