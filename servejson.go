package gosen

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/manvalls/gosen/buffersender"
	"github.com/manvalls/gosen/commands"
)

type JSONResponse struct {
	Commands []any  `json:"commands"`
	Version  string `json:"version,omitempty"`
}

func (h *handler) serveJSON(w http.ResponseWriter, r *http.Request) {
	buffer := buffersender.NewBufferSender()
	header := w.Header()
	wg := &sync.WaitGroup{}

	p := &Page{
		Version:    h.app.Version,
		Header:     header,
		StatusCode: http.StatusOK,
		Routine:    commands.NewRoutine(buffer, wg, nil),
	}

	h.f(p, r)
	wg.Wait()

	w.WriteHeader(p.StatusCode)
	data, _ := json.Marshal(JSONResponse{
		Commands: buffer.GetCommands(),
		Version:  p.Version,
	})

	w.Write(data)
}
