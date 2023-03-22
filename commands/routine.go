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
	Dynamic bool   `json:"dynamic,omitempty"`
}

func (r *Routine) runFork(sr *Routine, thread Thread) {
	defer r.waitGroup.Done()
	thread.Run(sr)
}

type Thread interface {
	Run(subroutine *Routine)
}

func (r *Routine) Fork(thread Thread) {
	r.waitGroup.Add(1)
	subroutine := r.subroutine()
	go r.runFork(subroutine, thread)
}

type funcThread struct {
	f func(subroutine *Routine)
}

func (f *funcThread) Run(subroutine *Routine) {
	f.f(subroutine)
}

func (r *Routine) ForkFunc(f func(subroutine *Routine)) {
	r.Fork(&funcThread{f})
}

func (r *Routine) Run(format string, args ...interface{}) {
	u := format
	dynamic := len(args) > 0

	if dynamic {
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
		r.runner.Run(r, u, dynamic)
		return
	}

	r.sender.SendCommand(RunCommand{u, r.id, dynamic})
}

type StartRoutineCommand struct {
	StartRoutine uint64 `json:"startRoutine"`
	Routine      uint64 `json:"routine,omitempty"`
}

func (r *Routine) subroutine() *Routine {
	nextId := r.getNextId()
	r.sender.SendCommand(StartRoutineCommand{nextId, r.id})
	return &Routine{r.sender, r.waitGroup, nextId, r.mux, r.nextId, r.runner}
}

func (r *Routine) Tx() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), false}
}

func (r *Routine) Once() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), true}
}
