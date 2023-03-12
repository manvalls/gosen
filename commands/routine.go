package commands

import (
	"sync"

	"github.com/zeebo/xxh3"
)

type Routine struct {
	sender CommandSender
	id     uint64
	mux    *sync.Mutex
	nextId *uint64

	runner *Runner
}

func (r *Routine) getNextId() uint64 {
	r.mux.Lock()
	defer r.mux.Unlock()
	*r.nextId++
	return *r.nextId
}

func NewRoutine(sender CommandSender, runner *Runner) *Routine {
	return &Routine{sender, 0, &sync.Mutex{}, new(uint64), runner}
}

type RunCommand struct {
	Run     string `json:"run"`
	Routine uint64 `json:"routine,omitempty"`
}

func (r *Routine) Run(url string) {
	if r.runner != nil {
		r.runner.Run(r, url)
		return
	}

	r.sender.SendCommand(RunCommand{url, r.id})
}

type StartRoutineCommand struct {
	StartRoutine uint64 `json:"startRoutine"`
	Routine      uint64 `json:"routine,omitempty"`
}

func (r *Routine) Subroutine() Routine {
	nextId := r.getNextId()
	r.sender.SendCommand(StartRoutineCommand{nextId, r.id})
	return Routine{r.sender, nextId, r.mux, r.nextId, r.runner}
}

func (r *Routine) Tx() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), false}
}

func (r *Routine) Once() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New(), true}
}
