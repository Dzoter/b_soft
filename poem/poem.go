package poem

import "pet/interfaces"

// Poem структура для хранения данных о стихах
type Poem struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var Poems []interfaces.Displayable
