package commands

import "sync"

type Routine struct {
	sender CommandSender
	id     uint
	mux    *sync.Mutex
	nextId *uint
}

func (r *Routine) getNextId() uint {
	r.mux.Lock()
	defer r.mux.Unlock()
	*r.nextId++
	return *r.nextId
}

func NewRoutine(sender CommandSender) *Routine {
	return &Routine{sender, 0, &sync.Mutex{}, new(uint)}
}

type RunCommand struct {
	Run     string `json:"run"`
	Routine uint   `json:"routine,omitempty"`
}

func (r *Routine) Run(url string) {
	r.sender.SendCommand(RunCommand{url, r.id})
}

type StartRoutineCommand struct {
	StartRoutine uint `json:"startRoutine"`
	Routine      uint `json:"routine,omitempty"`
}

func (r *Routine) Subroutine() Routine {
	nextId := r.getNextId()
	r.sender.SendCommand(StartRoutineCommand{nextId, r.id})
	return Routine{r.sender, nextId, r.mux, r.nextId}
}

func (r *Routine) Tx() *Transaction {
	return &Transaction{r.sender, nil, &sync.Mutex{}, 0, r.id}
}

func (r *Routine) UnmarshalJSON(data []byte) error {
	// TODO
	return nil
}
