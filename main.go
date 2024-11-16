package main

import (
	"fmt"
	"pet/handlers/poem"
	"pet/handlers/wiki"
	"pet/interfaces"
	"pet/terminal"
)

func main() {
	items := []interfaces.TitleDisplayable{poem.NewPoemFetcher(), wiki.NewWikiFetcher()}

	for {
		item, action := terminal.SelectItemsWithoutPaging(items)
		if item != nil {
			switch v := item.(type) {
			case poem.Fetcher:
				v.Process()
			case wiki.Fetcher:
				v.Process()
			default:
				fmt.Println("default")
			}
		}
		if action == terminal.ActionExit {
			terminal.DisplayMessage("Пока!")
			break
		}
	}
}
