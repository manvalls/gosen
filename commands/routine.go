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
}

func (r *Routine) getNextId() uint64 {
	r.mux.Lock()
	defer r.mux.Unlock()
	*r.nextId++
	return *r.nextId
}

func NewRoutine(sender CommandSender) *Routine {
	return &Routine{sender, 0, &sync.Mutex{}, new(uint64)}
}

type RunCommand struct {
	Run     string `json:"run"`
	Routine uint64 `json:"routine,omitempty"`
}

func (r *Routine) Run(url string) {
	r.sender.SendCommand(RunCommand{url, r.id})
}

type StartRoutineCommand struct {
	StartRoutine uint64 `json:"startRoutine"`
	Routine      uint64 `json:"routine,omitempty"`
}

func (r *Routine) Subroutine() Routine {
	nextId := r.getNextId()
	r.sender.SendCommand(StartRoutineCommand{nextId, r.id})
	return Routine{r.sender, nextId, r.mux, r.nextId}
}

func (r *Routine) Tx() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id, xxh3.New()}
}

func (r *Routine) UnmarshalJSON(data []byte) error {
	// TODO
	return nil
}
