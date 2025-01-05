package wiki

import (
	"b_soft/interfaces"
	"b_soft/terminal"
	"fmt"
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
	client, err := New("https://ru.wikipedia.org/w/api.php", "b soft")
	if err != nil {
		log.Fatal(err)
	}

	// Получаем список заголовков (пример верхнего уровня)
	titles, err := client.SearchTitles(input)
	if err != nil {
		log.Fatal(err)
	}

	if len(titles) == 0 {
		terminal.DisplayMessage("Заголовки не найдены")
		return
	}

	// Преобразуем []string в []interfaces.TitleDisplayable
	allTitles := make([]interfaces.TitleDisplayable, len(titles))
	for i, title := range titles {
		allTitles[i] = Title{Title: title}
	}

	chosenTitle, _ := terminal.SelectItemsWithoutPaging(allTitles)
	if chosenTitle != nil {
		switch v := chosenTitle.(type) {
		case Title:
			stringyPage, err := client.ReadTextOnly(v.Title)
			if err != nil {
				log.Fatal(err)
			}

			// Парсим текст в ContentItem
			contentItems := ConvertStringWikiToJSON(stringyPage)

			// Обрабатываем элементы контента рекурсивно
			handleContentItems(contentItems)
		default:
			fmt.Println("Выбранный элемент не является заголовком")
		}
	}
}
func handleContentItems(contentItems []ContentItem) {
	// Преобразуем []ContentItem в []interfaces.TitleDisplayable
	displayableItems := make([]interfaces.TitleDisplayable, len(contentItems))
	for i, item := range contentItems {
		displayableItems[i] = item
	}

	for {
		// Запрашиваем выбор пользователя
		chosableContent, action := terminal.SelectItemsWithoutPaging(displayableItems)
		if action == terminal.ActionExit {
			fmt.Println("Выход на предыдущий уровень")
			return // Прерываем текущую рекурсию и возвращаемся на уровень выше
		}

		if chosableContent != nil {
			// Приводим к ContentItem
			switch v := chosableContent.(type) {
			case ContentItem:
				fmt.Printf("Вы выбрали: %s (Type: %s)\n", v.Text, v.Type)

				// Если есть вложенные `Children`, заходим глубже
				if len(v.Children) > 0 {
					fmt.Println("Есть вложенные элементы, заходим глубже...")
					handleContentItems(v.Children) // Рекурсия
				} else {
					fmt.Println("Нет вложенных элементов")
				}
			default:
				fmt.Println("Выбранный элемент не является ContentItem")
			}
		} else {
			fmt.Println("Никакой элемент не выбран")
		}
	}
}
