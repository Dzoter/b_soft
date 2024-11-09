package structs

import "github.com/eiannone/keyboard"

type KeyBinding struct {
	Key         keyboard.Key // Клавиша для действия
	Char        rune         // Альтернативный символ клавиши (например, 'c')
	Description string       // Описание действия
	Action      func() bool  // Функция, выполняющая действие, возвращает true, если действие выполнено
}
