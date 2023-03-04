package commands

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash"
	"sync"

	"github.com/manvalls/gosen/template"
)

type Transaction struct {
	sender   CommandSender
	commands []interface{}
	mux      *sync.Mutex
	nextId   uint64
	routine  uint64
	hash     hash.Hash
	once     bool
}

type TransactionCommand struct {
	Transaction []any  `json:"tx"`
	Routine     uint64 `json:"routine,omitempty"`
	Hash        string `json:"hash,omitempty"`
	Once        bool   `json:"once,omitempty"`
}

func (t *Transaction) getNextId() uint64 {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.nextId++
	return t.nextId
}

func (t *Transaction) SendCommand(command any) {
	t.mux.Lock()
	defer t.mux.Unlock()

	switch c := command.(type) {

	case SelectorSubCommand:
		t.hash.Write([]byte{0})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		t.hash.Write([]byte(fmt.Sprintf(c.Selector, c.Args...)))

	case SelectorAllSubCommand:
		t.hash.Write([]byte{1})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		t.hash.Write([]byte(fmt.Sprintf(c.SelectorAll, c.Args...)))

	case FragmentSubCommand:
		t.hash.Write([]byte{2})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		c.Fragment.WriteHash(t.hash)

	case ContentSubCommand:
		t.hash.Write([]byte{3})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		binary.Write(t.hash, binary.LittleEndian, c.Content)

	case ParentSubCommand:
		t.hash.Write([]byte{4})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		binary.Write(t.hash, binary.LittleEndian, c.Target)

	case FirstChildSubCommand:
		t.hash.Write([]byte{5})
		binary.Write(t.hash, binary.LittleEndian, c.FirstChild)
		binary.Write(t.hash, binary.LittleEndian, c.Target)

	case LastChildSubCommand:
		t.hash.Write([]byte{6})
		binary.Write(t.hash, binary.LittleEndian, c.LastChild)
		binary.Write(t.hash, binary.LittleEndian, c.Target)

	case NextSiblingSubCommand:
		t.hash.Write([]byte{7})
		binary.Write(t.hash, binary.LittleEndian, c.NextSibling)
		binary.Write(t.hash, binary.LittleEndian, c.Target)

	case PrevSiblingSubCommand:
		t.hash.Write([]byte{8})
		binary.Write(t.hash, binary.LittleEndian, c.PrevSibling)
		binary.Write(t.hash, binary.LittleEndian, c.Target)

	case CloneSubCommand:
		t.hash.Write([]byte{9})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		binary.Write(t.hash, binary.LittleEndian, c.Clone)

	case TextSubCommand:
		t.hash.Write([]byte{10})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.Text))

	case HtmlSubCommand:
		t.hash.Write([]byte{11})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		c.Html.WriteHash(t.hash)

	case AttrSubCommand:
		t.hash.Write([]byte{12})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.Attr))
		t.hash.Write([]byte(c.Value))

	case RemoveAttrSubCommand:
		t.hash.Write([]byte{13})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.RemoveAttr))

	case RemoveSubCommand:
		t.hash.Write([]byte{14})
		binary.Write(t.hash, binary.LittleEndian, c.Remove)

	case InsertBeforeSubCommand:
		t.hash.Write([]byte{15})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		binary.Write(t.hash, binary.LittleEndian, c.Ref)

		switch ib := c.InsertBefore.(type) {
		case template.Template:
			ib.WriteHash(t.hash)

		case uint64:
			binary.Write(t.hash, binary.LittleEndian, ib)
		}

	case AppendSubCommand:
		t.hash.Write([]byte{16})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)

		switch a := c.Append.(type) {
		case template.Template:
			a.WriteHash(t.hash)

		case uint64:
			binary.Write(t.hash, binary.LittleEndian, a)
		}

	case WaitSubCommand:
		t.hash.Write([]byte{17})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.Wait))
		binary.Write(t.hash, binary.LittleEndian, c.Timeout)

	case AddToAttrSubCommand:
		t.hash.Write([]byte{18})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.AddToAttr))
		t.hash.Write([]byte(c.Value))

	case RemoveFromAttrSubCommand:
		t.hash.Write([]byte{19})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.RemoveFromAttr))
		t.hash.Write([]byte(c.Value))

	}

	t.commands = append(t.commands, command)
}

func (t *Transaction) Commit() {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.commands) == 0 {
		return
	}

	commands := t.commands
	t.commands = nil
	t.sender.SendCommand(TransactionCommand{commands, t.routine, base64.StdEncoding.EncodeToString(t.hash.Sum(nil)), t.once})
}

// Selectors

type SelectorSubCommand struct {
	Id       uint64        `json:"id"`
	Selector string        `json:"selector"`
	Args     []interface{} `json:"args,omitempty"`
	Parent   uint64        `json:"parent,omitempty"`
}

func (t *Transaction) S(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.SendCommand(SelectorSubCommand{nextId, selector, args, 0})
	return Node{t, nextId}
}

type SelectorAllSubCommand struct {
	Id          uint64        `json:"id"`
	SelectorAll string        `json:"selectorAll"`
	Args        []interface{} `json:"args,omitempty"`
	Parent      uint64        `json:"parent,omitempty"`
}

func (t *Transaction) All(selector string, args ...interface{}) Node {
	nextId := t.getNextId()
	t.SendCommand(SelectorAllSubCommand{nextId, selector, args, 0})
	return Node{t, nextId}
}

// Creation

type FragmentSubCommand struct {
	Id       uint64            `json:"id"`
	Fragment template.Template `json:"fragment"`
}

func (t *Transaction) Fragment(fragment template.Template) Node {
	nextId := t.getNextId()
	t.SendCommand(FragmentSubCommand{nextId, fragment})
	return Node{t, nextId}
}
