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

type waitCommand struct {
	event string
}

type attrCommand struct {
	name  string
	value string
}

type rmAttrCommand struct {
	name string
}

type addToAttrCommand struct {
	name  string
	value string
}

type rmFromAttrCommand struct {
	name  string
	value string
}

type addClassCommand struct {
	class string
}

type rmClassCommand struct {
	class string
}

type removeCommand struct{}

type emptyCommand struct{}

type replaceWithCommand struct {
	elementId uint
}

type insertBeforeCommand struct {
	childId uint
	refId   uint
}

type insertAfterCommand struct {
	childId uint
	refId   uint
}

type appendCommand struct {
	childId uint
}

type prependCommand struct {
	childId uint
}
