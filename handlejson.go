package gosen

import (
	"sync"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/jsonsender"
)

func handleJSON(c *gosenContext) *Routine {
	sender := &jsonsender.JSONSender{
		Writer: c,
	}

	wg := &sync.WaitGroup{}

	c.pending.Add(1)
	go func() {
		<-c.done
		defer c.pending.Done()
		wg.Wait()
		sender.End()
	}()

	return commands.NewRoutine(sender, wg, nil)
}
