package main

import (
	"fmt"
	"pet/interfaces"
	"pet/poem"
	"pet/terminal"
)

type welcome struct {
	title string
}

func (w welcome) DisplayTitle() string {
	return w.title
}
func (w welcome) DisplayBody() string {
	return w.title
}
func (w welcome) Process() {
	fmt.Println("заглушка")
}

func main() {

	wikiWelcome := welcome{title: "test"}

	items := []interfaces.Displayable{poem.NewPoemParser(), wikiWelcome}

	for {
		item, action := terminal.SelectItemsWithoutPaging(items)
		if item != nil {
			switch v := item.(type) {
			case poem.Parser:
				v.Process()
			case welcome:
				v.Process()
			}
		}
		if action == terminal.ActionExit {
			terminal.DisplayMessage("Пока!")
			break
		}
	}
}
