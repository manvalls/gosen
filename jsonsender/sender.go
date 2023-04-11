package jsonsender

import (
	"encoding/json"
	"io"
	"sync"
)

type JSONSender struct {
	mux            sync.Mutex
	prefaceWritten bool
	Writer         io.Writer
}

func (s *JSONSender) writePreface() bool {
	if !s.prefaceWritten {
		s.prefaceWritten = true
		s.Writer.Write([]byte{'['})
		return true
	}

	return false
}

func (s *JSONSender) SendCommand(command any) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.writePreface() {
		s.Writer.Write([]byte(`,`))
	}

	result, _ := json.Marshal(command)
	s.Writer.Write(result)
}

func (s *JSONSender) End() {
	s.writePreface()
	s.Writer.Write([]byte{']'})
}
