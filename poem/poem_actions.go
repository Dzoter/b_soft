package poem

import (
	"fmt"
	"github.com/atotto/clipboard"
	"pet/structs"
)

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
