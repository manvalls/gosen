package htmlsender

import "github.com/manvalls/gosen/commands"

func (s *HTMLSender) transaction(c commands.TransactionCommand) {

	for _, command := range c.Transaction {

		switch cmd := command.(type) {
		case commands.SelectorSubCommand:
			// TODO
		case commands.SelectorAllSubCommand:
			// TODO
		case commands.FragmentSubCommand:
			// TODO
		case commands.ContentSubCommand:
			// TODO
		case commands.CloneSubCommand:
			// TODO
		case commands.TextSubCommand:
			// TODO
		case commands.HtmlSubCommand:
			// TODO
		case commands.AttrSubCommand:
			// TODO
		case commands.RmAttrSubCommand:
			// TODO
		case commands.AddToAttrSubCommand:
			// TODO
		case commands.RmFromAttrSubCommand:
			// TODO
		case commands.AddClassSubCommand:
			// TODO
		case commands.RmClassSubCommand:
			// TODO
		case commands.RemoveSubCommand:
			// TODO
		case commands.EmptySubCommand:
			// TODO
		case commands.ReplaceWithSubCommand:
			// TODO
		case commands.InsertBeforeSubCommand:
			// TODO
		case commands.InsertAfterSubCommand:
			// TODO
		case commands.AppendSubCommand:
			// TODO
		case commands.PrependSubCommand:
			// TODO
		}

	}

}
