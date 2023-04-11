package ssesender

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/manvalls/gosen/commands"
)

type SSESender struct {
	Mux     *sync.Mutex
	Writer  io.Writer
	Flusher http.Flusher
}

func (s *SSESender) SendCommand(command any) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	result, _ := json.Marshal(command)
	s.Writer.Write([]byte("event: command\ndata: " + string(result) + "\n\n"))

	if s.Flusher != nil {
		s.Flusher.Flush()
	}
}

func (s *SSESender) SendEvent(e commands.Event) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	if e.Event != "" {
		s.Writer.Write([]byte("event: " + strings.ReplaceAll(e.Event, "\n", "") + "\n"))
	}

	if e.Data != "" {
		s.Writer.Write([]byte("data: " + strings.ReplaceAll(e.Data, "\n", "") + "\n"))
	}

	if e.Id != "" {
		s.Writer.Write([]byte("id: " + strings.ReplaceAll(e.Id, "\n", "") + "\n"))
	}

	if e.Retry != 0 {
		s.Writer.Write([]byte("retry: " + strconv.Itoa(e.Retry) + "\n"))
	}

	s.Writer.Write([]byte("\n"))

	if flusher, ok := s.Writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
