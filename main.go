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
			finder.FindPoem(input)
			if len(finder.Poems) == 0 {
				fmt.Println("Стихи не найдены.")
				continue
			}
			selectedPoem := terminal.SelectPoem(finder.Poems)
			if selectedPoem == nil {
				continue
			}
		}
	}
}
