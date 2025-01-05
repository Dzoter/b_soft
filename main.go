package main

import (
	"b_soft/handlers/poem"
	"b_soft/handlers/wiki"
	"b_soft/interfaces"
	"b_soft/terminal"
	"fmt"
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
