package wiki

import (
	"log"
	"pet/interfaces"
	"pet/terminal"
)

func (w Fetcher) Process() {
	//for {
	input := terminal.ReadInput()

	switch input {
	case terminal.ActionExit:
		return
	case "":
		return
	default:
		processInput(input)
	}
	//}
}

func processInput(input string) {
	// CREATE A NEW API STRUCT
	client, err := New("https://ru.wikipedia.org/w/api.php", "b soft")
	if err != nil {
		log.Fatal(err)
	}
	// получаем [] string тайтлов
	titles, err := client.SearchTitles(input)
	if err != nil {
		log.Fatal(err)
	}

	if len(titles) == 0 {
		terminal.DisplayMessage("Заголовки не найдены")
		return
	}
	// Преобразуем []string в []interfaces.Displayable
	allTitles := make([]interfaces.Displayable, 0, len(titles))
	for _, title := range titles {
		tmpTitle := Title{Title: title}
		allTitles = append(allTitles, tmpTitle)
	}

	item, _ := terminal.SelectItemsWithPaging(allTitles)

	if item != nil {
		switch v := item.(type) {
		case Title:
			_, err := client.ReadTextOnly(v.Title)
			if err != nil {
				log.Fatal(err)
			}
			test()
		default:
			break
		}
	}
}
