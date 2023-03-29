package jsonsender

import (
	"encoding/json"
	"io"
	"sync"
)

type JSONSender struct {
	mux            sync.Mutex
	prefaceWritten bool
	Writter        io.Writer
}

func (s *JSONSender) writePreface() bool {
	if !s.prefaceWritten {
		s.prefaceWritten = true
		s.Writter.Write([]byte{'['})
		return true
	}

	return false
}

func (s *JSONSender) SendCommand(command any) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.writePreface() {
		s.Writter.Write([]byte(`,`))
	}

	result, _ := json.Marshal(command)
	s.Writter.Write(result)
}

func (s *JSONSender) End() {
	s.writePreface()
	s.Writter.Write([]byte{']'})
}
