package buffersender

import "sync"

type BufferSender struct {
	mux      sync.Mutex
	commands []any
}

func New() *BufferSender {
	return &BufferSender{}
}

func (s *BufferSender) SendCommand(command any) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.commands = append(s.commands, command)
}

func (s *BufferSender) GetCommands() []any {
	s.mux.Lock()
	defer s.mux.Unlock()
	commands := s.commands
	s.commands = nil
	return commands
}
