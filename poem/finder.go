package poem

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"net/http"
	"net/url"
	"pet/interfaces"
	"pet/structs"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//TODO больно много всякой херни в классе отвечающем за поиск стихов. Декомпозировать.

// Poem структура для хранения данных о стихах
type Poem struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

func (p Poem) DisplayTitle() string {
	return fmt.Sprintf("%s %s\n", p.Author, p.Title)
}
func (p Poem) DisplayBody() string {
	return p.Text
}

// GetActions возвращает список действий, связанных с Poem
func (p Poem) GetActions() []structs.KeyBinding {
	return []structs.KeyBinding{
		{
			Char:        'c',
			Description: "Скопировать текст в буфер обмена",
			Action:      p.copyToClipboard,
		},
	}
}

// Функция для копирования текста в буфер обмена
func (p Poem) copyToClipboard() bool {
	err := clipboard.WriteAll(p.Text)
	if err != nil {
		fmt.Println("Ошибка при копировании текста в буфер обмена:", err)
		return false
	}
	fmt.Println("Текст успешно скопирован в буфер обмена.")
	return true
}

var Poems []interfaces.Displayable

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

func FindPoem(query string, page int) {
	Poems = Poems[:0]
	fmt.Printf("Ищем стих: %s страница %d", query, page)

	// Список сайтов для поиска
	baseURL := "https://www.culture.ru/literature/poems?"

	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	params.Add("query", query)

	fullURL := baseURL + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Ошибка при запросе к %s: %v\n", fullURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при парсинге HTML от %s: %v\n", fullURL, err)
	}

	// Найдем элемент <script> с данными JSON
	scriptContent := ""
	doc.Find("script#__NEXT_DATA__").Each(func(i int, s *goquery.Selection) {
		if scriptType, exists := s.Attr("type"); exists && scriptType == "application/json" {
			scriptContent = s.Text()
		}
	})

	if scriptContent == "" {
		fmt.Printf("Не удалось найти нужный элемент <script> на сайте %s\n", fullURL)
	}

	// Парсинг JSON данных
	var jsonData map[string]interface{}
	if err = json.Unmarshal([]byte(scriptContent), &jsonData); err != nil {
		fmt.Printf("Ошибка при парсинге JSON данных с сайта %s: %v\n", fullURL, err)
	}

	props, ok := safeMapInterface(jsonData["props"])
	if !ok {
		fmt.Printf("Отсутствует ключ 'props' в данных с сайта %s\n", fullURL)
	}

	pageProps, ok := safeMapInterface(props["pageProps"])
	if !ok {
		fmt.Printf("Отсутствует ключ 'pageProps' в данных с сайта %s\n", fullURL)
	}

	poems, ok := pageProps["poems"].([]interface{})
	if !ok {
		fmt.Printf("Отсутствует ключ 'poems' в данных с сайта %s\n", fullURL)
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
			fmt.Printf("Ошибка при приведении типов для стиха на сайте %s\n", fullURL)
			continue
		}
		author, authorOk := safeString(authorMap, "title")
		if !authorOk {
			fmt.Printf("Ошибка при приведении типов для стиха на сайте %s\n", fullURL)
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
