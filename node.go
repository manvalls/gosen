package gosen

type Node struct {
	tx *Transaction
	id uint
}

// Selectors

func (n Node) S(selector string, args ...interface{}) Node {
	nextId := n.tx.getNextId()
	n.tx.addCommand(selectorCommand{n.id, nextId, selector, args})
	return Node{n.tx, nextId}
}

func (n Node) All(selector string, args ...interface{}) Node {
	nextId := n.tx.getNextId()
	n.tx.addCommand(selectorAllCommand{n.id, nextId, selector, args})
	return Node{n.tx, nextId}
}

func (n Node) Content() Node {
	nextId := n.tx.getNextId()
	n.tx.addCommand(contentCommand{n.id, nextId})
	return Node{n.tx, nextId}
}

// Creation

func (n Node) Clone() Node {
	nextId := n.tx.getNextId()
	n.tx.addCommand(cloneCommand{n.id, nextId})
	return Node{n.tx, nextId}
}

// InnerText and InnerHTML

func (n Node) Text(text string) Node {
	n.tx.addCommand(textCommand{n.id, text})
	return n
}

func (n Node) HTML(html Template) Node {
	n.tx.addCommand(htmlCommand{n.id, html})
	return n
}

// Attributes

func (n Node) Attr(name string, value string) Node {
	n.tx.addCommand(attrCommand{n.id, name, value})
	return n
}

func (n Node) RmAttr(name string) Node {
	n.tx.addCommand(rmAttrCommand{n.id, name})
	return n
}

func (n Node) AddToAttr(name string, value string) Node {
	n.tx.addCommand(addToAttrCommand{n.id, name, value})
	return n
}

func (n Node) RmFromAttr(name string, value string) Node {
	n.tx.addCommand(rmFromAttrCommand{n.id, name, value})
	return n
}

func (n Node) AddClass(name string) Node {
	n.tx.addCommand(addClassCommand{n.id, name})
	return n
}

func (n Node) RmClass(name string) Node {
	n.tx.addCommand(rmClassCommand{n.id, name})
	return n
}

// Node manipulation

func (n Node) Remove() {
	n.tx.addCommand(removeCommand{n.id})
}

func (n Node) Empty() Node {
	n.tx.addCommand(emptyCommand{n.id})
	return n
}

func (n Node) ReplaceWith(otherNode Node) Node {
	if n.tx != otherNode.tx {
		panic("Cannot replace a node with a node from another transaction")
	}

	if n.id != otherNode.id {
		n.tx.addCommand(replaceWithCommand{n.id, otherNode.id})
	}

	return n
}

func (n Node) InsertBefore(child Node, ref Node) Node {
	if n.tx != child.tx {
		panic("Cannot insert a node from another transaction")
	}

	if n.id == child.id {
		panic("Cannot insert a node before itself")
	}

	n.tx.addCommand(insertBeforeCommand{n.id, child.id, ref.id})
	return n
}

func (n Node) InsertAfter(child Node, ref Node) Node {
	if n.tx != child.tx {
		panic("Cannot insert a node from another transaction")
	}

	if n.id == child.id {
		panic("Cannot insert a node after itself")
	}

	n.tx.addCommand(insertAfterCommand{n.id, child.id, ref.id})
	return n
}

func (n Node) Append(child Node) Node {
	if n.tx != child.tx {
		panic("Cannot append a node from another transaction")
	}

	if n.id == child.id {
		panic("Cannot append a node to itself")
	}

	n.tx.addCommand(appendCommand{n.id, child.id})
	return n
}

func (n Node) Prepend(child Node) Node {
	if n.tx != child.tx {
		panic("Cannot prepend a node from another transaction")
	}

	if n.id == child.id {
		panic("Cannot prepend a node to itself")
	}

	n.tx.addCommand(prependCommand{n.id, child.id})
	return n
}

// Misc

func (n Node) Wait(event string, timeout uint) Node {
	n.tx.addCommand(waitCommand{n.id, event, timeout})
	return n
}
