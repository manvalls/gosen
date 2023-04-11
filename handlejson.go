package gosen

import (
	"sync"

	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/gosen/jsonsender"
)

func handleJSON(c *gosenContext) *Routine {
	sender := &jsonsender.JSONSender{
		Writter: c,
	}

	wg := &sync.WaitGroup{}

	c.pending.Add(1)
	go func() {
		defer c.pending.Done()
		wg.Wait()
		sender.End()
	}()

	return commands.NewRoutine(sender, wg, nil)
}
