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
		c.routine = handleJSON(c)
		return c.routine
	}

	if r.URL.Query().Get("format") == "sse" {
		c.Header().Set("Content-Type", "text/event-stream")
		c.routine = handleSSE(c, r)
		return c.routine
	}

	c.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.routine = handleHTML(c, r)
	return c.routine
}
