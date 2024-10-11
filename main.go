package main

import (
	"fmt"
	"pet/finder"
	"pet/terminal"
)

func main() {
	for {
		input := terminal.ReadInput()
		if input != "" {
			processInput(input)
		}
	}
}

func processInput(input string) {
	page := 1
	for {
		finder.FindPoem(input, page)
		if len(finder.Poems) == 0 {
			fmt.Println("Стихи не найдены.")
			return
		}

		action := terminal.SelectPoem(finder.Poems)
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
