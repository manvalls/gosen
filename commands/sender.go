package commands

type CommandSender interface {
	SendCommand(command any)
}
