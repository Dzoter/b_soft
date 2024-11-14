package poem

import (
	"pet/terminal"
)

const title = "Стихи"

type Parser struct {
	title string
}

//TODO выглядит как будто этот код должен находится в другом классе

// NewPoemParser Используется для main, для определения экземпляра класса
func NewPoemParser() Parser {
	return Parser{
		title: title,
	}
}
func (p Parser) DisplayTitle() string {
	return p.title
}
func (p Parser) DisplayBody() string {
	return ""
}

func (p Parser) Process() {
	for {
		input := terminal.ReadInput()
		switch input {
		case "":
			processInput(input)
		case terminal.ActionExit:
			return
		default:
			return
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
