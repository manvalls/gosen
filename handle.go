package gosen

import (
	"net/http"

	"github.com/manvalls/gosen/commands"
)

func Handle(r *http.Request) *Routine {
	v := r.Context().Value(commands.GosenRoutineKey)
	if v != nil {
		routine, ok := v.(*Routine)
		if ok {
			return routine
		}
	}

	v = r.Context().Value(gosenContextKey)
	if v == nil {
		panic("gosen middleware missing")
	}

	c := v.(*gosenContext)
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.routine != nil {
		return c.routine
	}

	if r.URL.Query().Get("format") == "json" {
		c.Header().Set("Content-Type", "application/json")
		return handleJSON(c)
	}

	if r.URL.Query().Get("format") == "sse" {
		c.Header().Set("Content-Type", "text/event-stream")
		return handleSSE(c, r)
	}

	c.Header().Set("Content-Type", "text/html")
	return handleHTML(c, r)
}
