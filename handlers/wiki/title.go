package wiki

import "pet/interfaces"

// Title структура для хранения данных о стихах
type Title struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var Titles []interfaces.Displayable
