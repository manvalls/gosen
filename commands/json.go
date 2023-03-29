package commands

import (
	"github.com/manvalls/gosen/template"
	"github.com/valyala/fastjson"
)

func getSubCommand(v *fastjson.Value) any {

	if v.Exists("selector") {
		selector := v.GetStringBytes("selector")
		if selector == nil {
			return nil
		}

		return SelectorSubCommand{
			Id:       v.GetUint64("id"),
			Selector: string(selector),
			Dynamic:  v.GetBool("dynamic"),
			Parent:   v.GetUint64("parent"),
		}
	}

	if v.Exists("selectorAll") {
		selectorAll := v.GetStringBytes("selectorAll")
		if selectorAll == nil {
			return nil
		}

		return SelectorAllSubCommand{
			Id:          v.GetUint64("id"),
			SelectorAll: string(selectorAll),
			Dynamic:     v.GetBool("dynamic"),
			Parent:      v.GetUint64("parent"),
		}
	}

	if v.Exists("fragment") {
		fragment := v.GetStringBytes("fragment")
		if fragment == nil {
			return nil
		}

		return FragmentSubCommand{
			Id:       v.GetUint64("id"),
			Fragment: template.String(string(fragment)),
		}
	}

	if v.Exists("content") {
		return ContentSubCommand{
			Id:      v.GetUint64("id"),
			Content: v.GetUint64("content"),
		}
	}

	if v.Exists("target") {

		if v.Exists("parent") {
			return ParentSubCommand{
				Parent: v.GetUint64("parent"),
				Target: v.GetUint64("target"),
			}
		}

		if v.Exists("firstChild") {
			return FirstChildSubCommand{
				FirstChild: v.GetUint64("firstChild"),
				Target:     v.GetUint64("target"),
			}
		}

		if v.Exists("lastChild") {
			return LastChildSubCommand{
				LastChild: v.GetUint64("lastChild"),
				Target:    v.GetUint64("target"),
			}
		}

		if v.Exists("nextSibling") {
			return NextSiblingSubCommand{
				NextSibling: v.GetUint64("nextSibling"),
				Target:      v.GetUint64("target"),
			}
		}

		if v.Exists("prevSibling") {
			return PrevSiblingSubCommand{
				PrevSibling: v.GetUint64("prevSibling"),
				Target:      v.GetUint64("target"),
			}
		}

		if v.Exists("text") {
			text := v.GetStringBytes("text")
			if text == nil {
				return nil
			}

			return TextSubCommand{
				Target: v.GetUint64("target"),
				Text:   string(text),
			}
		}

		if v.Exists("html") {
			html := v.GetStringBytes("html")
			if html == nil {
				return nil
			}

			return HtmlSubCommand{
				Target: v.GetUint64("target"),
				Html:   template.String(string(html)),
			}
		}

		if v.Exists("attr") {
			attr := v.GetStringBytes("attr")
			if attr == nil {
				return nil
			}

			value := v.GetStringBytes("value")
			if value == nil {
				return nil
			}

			return AttrSubCommand{
				Target: v.GetUint64("target"),
				Attr:   string(attr),
				Value:  string(value),
			}
		}

		if v.Exists("removeAttr") {
			removeAttr := v.GetStringBytes("removeAttr")
			if removeAttr == nil {
				return nil
			}

			return RemoveAttrSubCommand{
				Target:     v.GetUint64("target"),
				RemoveAttr: string(removeAttr),
			}
		}

		if v.Exists("addToAttr") {
			addToAttr := v.GetStringBytes("addToAttr")
			if addToAttr == nil {
				return nil
			}

			value := v.GetStringBytes("value")
			if value == nil {
				return nil
			}

			return AddToAttrSubCommand{
				Target:    v.GetUint64("target"),
				AddToAttr: string(addToAttr),
				Value:     string(value),
			}
		}

		if v.Exists("removeFromAttr") {
			removeFromAttr := v.GetStringBytes("removeFromAttr")
			if removeFromAttr == nil {
				return nil
			}

			value := v.GetStringBytes("value")
			if value == nil {
				return nil
			}

			return RemoveFromAttrSubCommand{
				Target:         v.GetUint64("target"),
				RemoveFromAttr: string(removeFromAttr),
				Value:          string(value),
			}
		}

		if v.Exists("wait") {
			wait := v.GetStringBytes("wait")
			if wait == nil {
				return nil
			}

			return WaitSubCommand{
				Target:  v.GetUint64("target"),
				Wait:    string(wait),
				Timeout: v.GetUint64("timeout"),
			}
		}
	}

	if v.Exists("clone") {
		return CloneSubCommand{
			Id:    v.GetUint64("id"),
			Clone: v.GetUint64("clone"),
		}
	}

	if v.Exists("remove") {
		return RemoveSubCommand{
			Remove: v.GetUint64("remove"),
		}
	}

	if v.Exists("insertNodeBefore") {
		return InsertNodeBeforeSubCommand{
			InsertNodeBefore: v.GetUint64("insertNodeBefore"),
			Parent:           v.GetUint64("parent"),
			Ref:              v.GetUint64("ref"),
		}
	}

	if v.Exists("insertBefore") {
		inserBefore := v.GetStringBytes("insertBefore")
		if inserBefore == nil {
			return nil
		}

		return InsertBeforeSubCommand{
			InsertBefore: template.String(string(inserBefore)),
			Parent:       v.GetUint64("parent"),
			Ref:          v.GetUint64("ref"),
		}
	}

	if v.Exists("appendNode") {
		return AppendNodeSubCommand{
			AppendNode: v.GetUint64("appendNode"),
			Parent:     v.GetUint64("parent"),
		}
	}

	if v.Exists("append") {
		append := v.GetStringBytes("append")
		if append == nil {
			return nil
		}

		return AppendSubCommand{
			Append: template.String(string(append)),
			Parent: v.GetUint64("parent"),
		}
	}

	return nil
}

func (r *Routine) UnmarshalJSON(data []byte) error {
	var p fastjson.Parser

	routinesMap := make(map[uint64]uint64)
	getRoutine := func(id uint64) uint64 {
		if id == 0 {
			return r.id
		}

		if _, ok := routinesMap[id]; !ok {
			routinesMap[id] = r.getNextId()
		}

		return routinesMap[id]
	}

	v, err := p.Parse(string(data))
	if err != nil {
		return err
	}

	arr := v.GetArray()
	for _, val := range arr {
		if val.Exists("hash") || val.Exists("tx") {
			transaction := []any{}
			txArr := val.GetArray("tx")

			for _, txVal := range txArr {
				sc := getSubCommand(txVal)
				if sc != nil {
					transaction = append(transaction, sc)
				}
			}

			hash := val.GetStringBytes("hash")
			if hash == nil {
				hash = []byte{}
			}

			r.sender.SendCommand(TransactionCommand{
				Transaction: transaction,
				Routine:     getRoutine(val.GetUint64("routine")),
				Hash:        string(hash),
				Once:        val.GetBool("once"),
			})

			continue
		}

		if val.Exists("startRoutine") {
			r.sender.SendCommand(StartRoutineCommand{
				StartRoutine: val.GetUint64("startRoutine"),
				Routine:      getRoutine(val.GetUint64("routine")),
			})

			continue
		}

		if val.Exists("run") {
			run := val.GetStringBytes("run")
			if run == nil {
				run = []byte{}
			}

			r.sender.SendCommand(RunCommand{
				Routine: getRoutine(val.GetUint64("routine")),
				Run:     string(run),
			})

			continue
		}
	}

	return nil
}
