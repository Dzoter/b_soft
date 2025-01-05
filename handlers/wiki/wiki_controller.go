package wiki

import (
	"b_soft/interfaces"
	"b_soft/terminal"
	"log"
)

func (w Fetcher) Process() {
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
	allTitles := make([]interfaces.TitleDisplayable, 0, len(titles))
	for _, title := range titles {
		tmpTitle := Title{Title: title}
		allTitles = append(allTitles, tmpTitle)
	}

	chosenTitle, _ := terminal.SelectItemsWithoutPaging(allTitles)

	if chosenTitle != nil {
		switch v := chosenTitle.(type) {
		case Title:
			stringyPage, err := client.ReadTextOnly(v.Title)
			if err != nil {
				log.Fatal(err)
			}
			ConvertStringWikiToJSON(stringyPage)
		default:
			break
		}
	}
}
