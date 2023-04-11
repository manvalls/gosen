package gosen

import (
	"net/http"
	"sync"
	"time"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/ssesender"
)

func handleSSE(c *gosenContext, r *http.Request) *Routine {
	flusher, _ := c.ResponseWriter.(http.Flusher)

	if flusher != nil && c.config.SSEKeepAlive > 0 {
		go func() {
			for {
				select {
				case <-r.Context().Done():
					return
				case <-time.After(time.Duration(c.config.SSEKeepAlive) * time.Second):
					c.Write([]byte(":keepalive\n\n"))
					flusher.Flush()
				}
			}
		}()
	}

	mux := &sync.Mutex{}

	sender := &ssesender.SSESender{
		Mux:     mux,
		Writer:  c,
		Flusher: flusher,
	}

	return commands.NewRoutine(sender, &c.pending, nil)
}
