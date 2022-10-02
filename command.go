package gosen

type commandSender interface {
	sendCommand(elementId uint, command interface{})
}

type transactionCommand struct {
	nextId uint
}

type commitCommand struct{}

type selectorCommand struct {
	id       uint
	selector string
	args     []interface{}
}

type selectorAllCommand struct {
	id       uint
	selector string
	args     []interface{}
}

type contentCommand struct {
	id uint
}

type createCommand struct {
	id       uint
	template Template
}

type fragmentCommand struct {
	id       uint
	template Template
}

type cloneCommand struct {
	id uint
}

type textCommand struct {
	text string
}

type htmlCommand struct {
	html Template
}

type loadScriptCommand struct {
	scriptURL string
	async     bool
}
