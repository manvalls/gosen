package ssesender

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/manvalls/gosen/commands"
)

type SSESender struct {
	mux            sync.Mutex
	versionWritten bool
	RunList        []string
	Writter        io.Writer
	Flusher        http.Flusher
	commands.VersionGetter
}

func (s *SSESender) SendCommand(command any) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.versionWritten {
		s.Writter.Write([]byte("event: version\ndata: " + s.Version() + "\n\n"))
		s.versionWritten = true
	}

	if command, ok := command.(commands.RunCommand); ok {
		s.RunList = append(s.RunList, command.Run)
	}

	result, _ := json.Marshal(command)
	s.Writter.Write([]byte("event: command\ndata: " + string(result) + "\n\n"))
	s.Flusher.Flush()
}
