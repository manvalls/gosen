package gosen

type commandSender interface {
	sendCommand(command interface{})
}

type transactionCommand struct {
	commands []interface{}
}

type scopeCommand struct {
	commands []interface{}
}

type rootSelectorCommand struct {
	id       uint
	selector string
	args     []interface{}
}

type rootSelectorAllCommand struct {
	id       uint
	selector string
	args     []interface{}
}

type selectorCommand struct {
	parentID uint
	id       uint
	selector string
	args     []interface{}
}

type selectorAllCommand struct {
	parentID uint
	id       uint
	selector string
	args     []interface{}
}

type contentCommand struct {
	parentID uint
	id       uint
}

type fragmentCommand struct {
	id       uint
	template Template
}

type cloneCommand struct {
	targetID uint
	id       uint
}

type textCommand struct {
	targetID uint
	text     string
}

type htmlCommand struct {
	targetID uint
	html     Template
}

type waitCommand struct {
	targetID uint
	event    string
	timeout  uint
}

type attrCommand struct {
	targetID uint
	name     string
	value    string
}

type rmAttrCommand struct {
	targetID uint
	name     string
}

type addToAttrCommand struct {
	targetID uint
	name     string
	value    string
}

type rmFromAttrCommand struct {
	targetID uint
	name     string
	value    string
}

type addClassCommand struct {
	targetID uint
	class    string
}

type rmClassCommand struct {
	targetID uint
	class    string
}

type removeCommand struct {
	targetID uint
}

type emptyCommand struct {
	targetID uint
}

type replaceWithCommand struct {
	targetID  uint
	elementID uint
}

type insertBeforeCommand struct {
	parentID uint
	childID  uint
	refID    uint
}

type insertAfterCommand struct {
	parentID uint
	childID  uint
	refID    uint
}

type appendCommand struct {
	parentID uint
	childID  uint
}

type prependCommand struct {
	parentID uint
	childID  uint
}

type runCommand struct {
	url string
}

type onceCommand struct {
	url string
}
