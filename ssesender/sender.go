package ssesender

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type SSESender struct {
	Mux     *sync.Mutex
	Writter io.Writer
	Flusher http.Flusher
}

func (s *SSESender) SendCommand(command any) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	result, _ := json.Marshal(command)
	s.Writter.Write([]byte("event: command\ndata: " + string(result) + "\n\n"))

	if s.Flusher != nil {
		s.Flusher.Flush()
	}
}
