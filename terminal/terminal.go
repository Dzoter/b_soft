package terminal

import (
	"bufio"
	"fmt"
	"os"
	"pet/finder"
	"strings"

	"github.com/atotto/clipboard"
)

// ReadInput считывает ввод пользователя
func ReadInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите текст стиха для поиска:")

	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return ""
	}

	input = strings.TrimSpace(input)
	return input
}

// DisplayPoems отображает список найденных стихов
func DisplayPoems(poems []finder.Poem) {
	fmt.Println("Найденные стихи:")
	for i, poem := range poems {
		fmt.Printf("%d %s - %s\n", i, poem.Title, poem.Author)
	}
}

// SelectPoem предлагает пользователю выбрать стихотворение из списка
func SelectPoem(poems []finder.Poem) *finder.Poem {
	for {
		DisplayPoems(poems) // Отображаем список стихотворений на каждом повторе
		var choice int
		fmt.Print("Введите номер нужного стихотворения (или -1 для возврата к вводу поиска): ")
		_, err := fmt.Scan(&choice)
		if err != nil || choice < -1 || choice >= len(poems) {
			fmt.Println("Неверный выбор, попробуйте снова.")
			continue
		}

		if choice == -1 {
			poems = poems[:0]
			return nil
		}

		poem := &poems[choice]
		proceed := displayPoemOptions(poem)
		if proceed {
			return poem
		} else {
			continue // Возврат к списку стихов
		}
	}
}

// displayPoemOptions предоставляет пользователю варианты действий с выбранным стихотворением
func displayPoemOptions(poem *finder.Poem) bool {
	fmt.Printf("Вы выбрали: %s - %s\nТекст:\n%s\n", poem.Title, poem.Author, poem.Text)
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Скопировать текст в буфер обмена")
		fmt.Println("0 - Вернуться к списку стихотворений")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении ввода:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "1" {
			copyToClipboard(poem.Text)
			return true
		} else if input == "0" {
			return false
		} else {
			fmt.Println("Неверный выбор, попробуйте снова.")
		}
	}
}

// copyToClipboard копирует текст в буфер обмена
func copyToClipboard(text string) {
	err := clipboard.WriteAll(text)
	if err != nil {
		fmt.Println("Ошибка при копировании текста в буфер обмена:", err)
		return
	}
	fmt.Println("Текст успешно скопирован в буфер обмена.")
}
