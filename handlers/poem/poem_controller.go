package poem

import (
	"pet/terminal"
)

func (p Fetcher) Process() {
	for {
		input := terminal.ReadInput()
		switch input {
		case terminal.ActionExit:
			return
		case "":
			return
		default:
			processInput(input)
		}
	}
}

func processInput(input string) {
	page := 1
	for {
		FindPoem(input, page) // Используем функцию FindPoem из второго пакета
		if len(Poems) == 0 {
			terminal.DisplayMessage("Стихи не найдены")
			return
		}

		item, action := terminal.SelectItemsWithPaging(Poems)
		if item != nil {
			switch v := item.(type) {
			case Poem:
				terminal.DisplayItemBody(v)
				terminal.DisplayActions(v)
			default:
				break
			}
		}
		switch action {
		case terminal.ActionNext:
			page++
		case terminal.ActionPrev:
			if page > 1 {
				page--
			}
		case terminal.ActionExit:
			return
		}

	}
}
