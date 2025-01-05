package wiki

import "b_soft/interfaces"

// Title структура для хранения данных о стихах
type Title struct {
	Title string `json:"title"`
}

var Titles []interfaces.Displayable
