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

func handleHTML(c *gosenContext, r *http.Request) *Routine {
	var html *htmlsender.HTMLSender
	var buffer *buffersender.BufferSender
	var sender commands.CommandSender

	html = htmlsender.NewHTMLSender(c.selectorCache)

	if c.config.Hydrate {
		buffer = buffersender.NewBufferSender()
		sender = multisender.NewMultiSender(buffer, html)
	} else {
		sender = html
	}

	wg := &sync.WaitGroup{}

	runner := &commands.Runner{
		BaseRequest: r,
		MapRunURL:   c.config.MapRunURL,
		Handler:     c.handler,
	}

	c.pending.Add(1)
	go func() {
		<-c.done
		defer c.pending.Done()
		wg.Wait()

		if c.config.Hydrate {
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
				if c.version != "" {
					v, _ := json.Marshal(c.version)
					script += "window.__GOSEN_PAGE_VERSION__=" + string(v) + ";"
				}

				html.PrependScript(script)
			}
		}

		html.Render(c)
	}()

	return commands.NewRoutine(sender, wg, runner)
}
