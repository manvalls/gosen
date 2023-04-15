package commands

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/zeebo/xxh3"
)

type Routine struct {
	sender    CommandSender
	waitGroup *sync.WaitGroup
	id        uint64
	mux       *sync.Mutex
	nextId    *uint64

	runner *Runner
}

func (r *Routine) getNextId() uint64 {
	r.mux.Lock()
	defer r.mux.Unlock()
	*r.nextId++
	return *r.nextId
}

func NewRoutine(sender CommandSender, wg *sync.WaitGroup, runner *Runner) *Routine {
	return &Routine{sender, wg, 0, &sync.Mutex{}, new(uint64), runner}
}

type RunCommand struct {
	Run     string `json:"run"`
	Routine uint64 `json:"routine,omitempty"`
}

func (r *Routine) Run(format string, args ...interface{}) {
	u := format

	if len(args) > 0 {
		escapedArgs := make([]interface{}, len(args))
		for i, arg := range args {
			if str, ok := arg.(string); ok {
				escapedArgs[i] = url.QueryEscape(str)
			} else {
				escapedArgs[i] = arg
			}
		}

		u = fmt.Sprintf(format, escapedArgs...)
	}

	if r.runner != nil {
		r.runner.Run(r, u)
		return
	}

	r.sender.SendCommand(RunCommand{u, r.id})
}

type StartRoutineCommand struct {
	StartRoutine uint64 `json:"startRoutine"`
	Routine      uint64 `json:"routine,omitempty"`
}

type EndRoutineCommand struct {
	EndRoutine uint64 `json:"endRoutine"`
}

func (r *Routine) subroutine() *Routine {
	nextId := r.getNextId()
	r.sender.SendCommand(StartRoutineCommand{nextId, r.id})
	return &Routine{r.sender, r.waitGroup, nextId, r.mux, r.nextId, r.runner}
}

func runSubroutine(subroutine *Routine, f func(subroutine *Routine)) {
	f(subroutine)
	subroutine.sender.SendCommand(EndRoutineCommand{subroutine.id})
	subroutine.waitGroup.Done()
}

func (r *Routine) Fork(f func(subroutine *Routine)) {
	r.waitGroup.Add(1)
	subroutine := r.subroutine()
	go runSubroutine(subroutine, f)
}

func (r *Routine) Tx() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), false}
}

func (r *Routine) Once() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), true}
}

func (r *Routine) SendEvent(e Event) {
	evenrSender, ok := r.sender.(EventSender)
	if ok {
		evenrSender.SendEvent(e)
	}
}
