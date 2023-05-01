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
	commands []any
	mux      *sync.Mutex
	nextId   uint64
	routine  uint64
	hash     hash.Hash
	once     bool
}

type TransactionCommand struct {
	Transaction []any  `json:"tx,omitempty"`
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
		t.hash.Write([]byte(c.Selector))

	case SelectorAllSubCommand:
		t.hash.Write([]byte{1})
		binary.Write(t.hash, binary.LittleEndian, c.Id)
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		t.hash.Write([]byte(c.SelectorAll))

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

	case InsertNodeBeforeSubCommand:
		t.hash.Write([]byte{15})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		binary.Write(t.hash, binary.LittleEndian, c.Ref)
		binary.Write(t.hash, binary.LittleEndian, c.InsertNodeBefore)

	case InsertBeforeSubCommand:
		t.hash.Write([]byte{16})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		binary.Write(t.hash, binary.LittleEndian, c.Ref)
		c.InsertBefore.WriteHash(t.hash)

	case AppendNodeSubCommand:
		t.hash.Write([]byte{17})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		binary.Write(t.hash, binary.LittleEndian, c.AppendNode)

	case AppendSubCommand:
		t.hash.Write([]byte{18})
		binary.Write(t.hash, binary.LittleEndian, c.Parent)
		c.Append.WriteHash(t.hash)

	case WaitSubCommand:
		t.hash.Write([]byte{19})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.Wait))
		binary.Write(t.hash, binary.LittleEndian, c.Timeout)

	case AddToAttrSubCommand:
		t.hash.Write([]byte{20})
		binary.Write(t.hash, binary.LittleEndian, c.Target)
		t.hash.Write([]byte(c.AddToAttr))
		t.hash.Write([]byte(c.Value))

	case RemoveFromAttrSubCommand:
		t.hash.Write([]byte{21})
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
	t.sender.SendCommand(TransactionCommand{commands, t.routine, base64.RawStdEncoding.EncodeToString(t.hash.Sum(nil)), t.once})
}

// Selectors

type SelectorSubCommand struct {
	Id       uint64 `json:"id"`
	Selector string `json:"selector"`
	Dynamic  bool   `json:"dynamic,omitempty"`
	Parent   uint64 `json:"parent,omitempty"`
}

func (t *Transaction) S(selector string, args ...any) Node {
	nextId := t.getNextId()

	if len(args) > 0 {
		t.SendCommand(SelectorSubCommand{nextId, fmt.Sprintf(selector, args...), true, 0})
		return Node{t, nextId}
	}

	t.SendCommand(SelectorSubCommand{nextId, selector, false, 0})
	return Node{t, nextId}
}

type SelectorAllSubCommand struct {
	Id          uint64 `json:"id"`
	SelectorAll string `json:"selectorAll"`
	Dynamic     bool   `json:"dynamic,omitempty"`
	Parent      uint64 `json:"parent,omitempty"`
}

func (t *Transaction) All(selector string, args ...any) Node {
	nextId := t.getNextId()

	if len(args) > 0 {
		t.SendCommand(SelectorAllSubCommand{nextId, fmt.Sprintf(selector, args...), true, 0})
		return Node{t, nextId}
	}

	t.SendCommand(SelectorAllSubCommand{nextId, selector, false, 0})
	return Node{t, nextId}
}

type IdSubCommand struct {
	Id        uint64 `json:"id"`
	ElementId string `json:"elementId"`
}

func (t *Transaction) Id(id string, args ...any) Node {
	nextId := t.getNextId()

	if len(args) > 0 {
		t.SendCommand(IdSubCommand{nextId, fmt.Sprintf(id, args...)})
		return Node{t, nextId}
	}

	t.SendCommand(IdSubCommand{nextId, id})
	return Node{t, nextId}
}

type HeadSubCommand struct {
	Head uint64 `json:"head"`
}

func (t *Transaction) Head() Node {
	nextId := t.getNextId()
	t.SendCommand(HeadSubCommand{nextId})
	return Node{t, nextId}
}

type BodySubCommand struct {
	Body uint64 `json:"body"`
}

func (t *Transaction) Body() Node {
	nextId := t.getNextId()
	t.SendCommand(BodySubCommand{nextId})
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
