package finder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Poem структура для хранения данных о стихах
type Poem struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var Poems []Poem

// safeString извлекает строку из карты, если это возможно.
func safeString(m map[string]interface{}, key string) (string, bool) {
	value, ok := m[key]
	if !ok {
		return "", false
	}
	strValue, ok := value.(string)
	return strValue, ok
}

// safeMap извлекает карту из карты, если это возможно.
func safeMap(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	value, ok := m[key]
	if !ok {
		return nil, false
	}
	mapValue, ok := value.(map[string]interface{})
	return mapValue, ok
}

// safeMapInterface извлекает map[string]interface{} из interface{}, если это возможно.
func safeMapInterface(data interface{}) (map[string]interface{}, bool) {
	mapValue, ok := data.(map[string]interface{})
	return mapValue, ok
}

func cleanText(text string) string {
	// Шаблон для поиска строк, содержащих теги в квадратных скобках или URL в круглых скобках
	reBracketsAndUrls := regexp.MustCompile(`(\[[^\]]+\]|\(http[^\)]+\))`)
	// Найти все совпадения
	matches := reBracketsAndUrls.FindAllStringIndex(text, -1)

	if len(matches) > 0 {
		// Обрезать текст до первого совпадения
		text = text[:matches[0][0]]
	}

	// Шаблон для поиска HTML-тегов
	reHTML := regexp.MustCompile(`<.*?>`)
	// Удалить все HTML-теги
	text = reHTML.ReplaceAllString(text, "")

	// Удалить лишние пробелы и пустые строки
	return strings.TrimSpace(text)
}

func FindPoem(query string) {
	Poems = Poems[:0]
	fmt.Println("Ищем стих:", query)

	// Список сайтов для поиска
	sites := []string{
		"https://www.culture.ru/literature/poems?query=",
	}

	escapedQuery := url.QueryEscape(query)

	for _, site := range sites {
		fullURL := site + escapedQuery
		resp, err := http.Get(fullURL)
		if err != nil {
			fmt.Printf("Ошибка при запросе к %s: %v\n", site, err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Printf("Ошибка при парсинге HTML от %s: %v\n", site, err)
			continue
		}

		// Найдем элемент <script> с данными JSON
		scriptContent := ""
		doc.Find("script#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
			if scriptType, exists := s.Attr("type"); exists && scriptType == "application/json" {
				scriptContent = s.Text()
			}
		})

		if scriptContent == "" {
			fmt.Printf("Не удалось найти нужный элемент <script> на сайте %s\n", site)
			continue
		}

		// Парсинг JSON данных
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(scriptContent), &jsonData); err != nil {
			fmt.Printf("Ошибка при парсинге JSON данных с сайта %s: %v\n", site, err)
			continue
		}

		props, ok := safeMapInterface(jsonData["props"])
		if !ok {
			fmt.Printf("Отсутствует ключ 'props' в данных с сайта %s\n", site)
			continue
		}

		pageProps, ok := safeMapInterface(props["pageProps"])
		if !ok {
			fmt.Printf("Отсутствует ключ 'pageProps' в данных с сайта %s\n", site)
			continue
		}

		poems, ok := pageProps["poems"].([]interface{})
		if !ok {
			fmt.Printf("Отсутствует ключ 'poems' в данных с сайта %s\n", site)
			continue
		}

		// Обработка данных о стихах
		for _, poemData := range poems {
			poemMap, ok := safeMapInterface(poemData)
			if !ok {
				continue
			}

			title, titleOk := safeString(poemMap, "title")
			text, textOk := safeString(poemMap, "text")
			authorMap, authorMapOk := safeMap(poemMap, "author")
			if !titleOk || !authorMapOk || !textOk {
				fmt.Printf("Ошибка при приведении типов для стиха на сайте %s\n", site)
				continue
			}
			author, authorOk := safeString(authorMap, "title")
			if !authorOk {
				fmt.Printf("Ошибка при приведении типов для стиха на сайте %s\n", site)
				continue
			}
			poem := Poem{
				Title:  title,
				Author: author,
				Text:   cleanText(text),
			}

			Poems = append(Poems, poem)
		}
	}
}
