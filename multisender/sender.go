package multisender

import "github.com/manvalls/gosen/commands"

type MultiSender struct {
	senders []commands.CommandSender
}

func NewMultiSender(senders ...commands.CommandSender) *MultiSender {
	return &MultiSender{senders}
}

func (m *MultiSender) SendCommand(command any) {
	for _, sender := range m.senders {
		sender.SendCommand(command)
	}
}
