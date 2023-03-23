package ssesender

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/manvalls/gosen/commands"
)

type SSESender struct {
	versionWritten bool
	Mux            *sync.Mutex
	RunList        []string
	Writter        io.Writer
	Flusher        http.Flusher
	commands.VersionGetter
}

func (s *SSESender) SendCommand(command any) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	if !s.versionWritten {
		s.Writter.Write([]byte("event: version\ndata: " + strings.ReplaceAll(s.Version(), "\n", "") + "\n\n"))
		s.versionWritten = true
	}

	if command, ok := command.(commands.RunCommand); ok {
		s.RunList = append(s.RunList, command.Run)
	}

	result, _ := json.Marshal(command)
	s.Writter.Write([]byte("event: command\ndata: " + string(result) + "\n\n"))

	if s.Flusher != nil {
		s.Flusher.Flush()
	}
}
