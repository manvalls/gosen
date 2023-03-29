package jsonsender

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/manvalls/gosen/commands"
)

type JSONSender struct {
	mux            sync.Mutex
	versionWritten bool
	Writter        io.Writer
	commands.VersionGetter
}

func (s *JSONSender) writeVersion() bool {
	if !s.versionWritten {
		s.versionWritten = true

		versionString := ""
		if s.Version() != "" {
			v, _ := json.Marshal(s.Version())
			versionString = `"version":` + string(v) + `,`
		}

		s.Writter.Write([]byte(`{` + versionString + `"commands":[`))
		return true
	}

	return false
}

func (s *JSONSender) SendCommand(command any) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if !s.writeVersion() {
		s.Writter.Write([]byte(`,`))
	}

	result, _ := json.Marshal(command)
	s.Writter.Write(result)
}

func (s *JSONSender) End() {
	s.writeVersion()
	s.Writter.Write([]byte(`]}`))
}
