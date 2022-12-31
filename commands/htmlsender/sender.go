package htmlsender

import (
	"github.com/manvalls/gosen/commands"
	"github.com/manvalls/mutexmap"
)

type HTMLSender struct {
	mutex mutexmap.MutexMap[uint]
}

func (s *HTMLSender) run(c commands.RunCommand) {
	// TODO
}

func (s *HTMLSender) runOnce(c commands.RunOnceCommand) {
	// TODO
}

func (s *HTMLSender) transaction(c commands.TransactionCommand) {
	// TODO
}

func (s *HTMLSender) SendCommand(command any) {
	switch c := command.(type) {

	case commands.RunCommand:
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)
		s.run(c)

	case commands.RunOnceCommand:
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)
		s.runOnce(c)

	case commands.StartRoutineCommand:
		s.mutex.Lock(c.StartRoutine)
		defer s.mutex.Unlock(c.StartRoutine)
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)

	case commands.TransactionCommand:
		s.mutex.Lock(c.Routine)
		defer s.mutex.Unlock(c.Routine)
		s.transaction(c)

	}
}
