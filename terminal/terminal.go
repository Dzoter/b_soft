package terminal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pet/finder"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/eiannone/keyboard"
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

// displayPoem отображает список найденных стихов
func displayPoem(poem finder.Poem) {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%s - %s\n", poem.Title, poem.Author)

}
func DisplayMessage(message string) {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("%s\n", message)
}

// SelectPoem предлагает пользователю выбрать стихотворение из списка
func SelectPoem(poems []finder.Poem) *finder.Poem {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	var choice int

	DisplayMessage("Стихотворения загружены, нажимайте вниз/вверх для выбора")

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return nil
		case keyboard.KeyArrowUp:
			if choice > 0 {
				choice--
			}
			displayPoem(poems[choice])
		case keyboard.KeyArrowDown:
			if choice < len(poems)-1 {
				choice++
			}
			displayPoem(poems[choice])
		case keyboard.KeyEnter:
			res := displayPoemWithOpt(&poems[choice])
			if res {
				break
			}
			displayPoem(poems[choice])
			continue
		default:
			continue
		}
	}
}

// displayPoemWithOpt предоставляет пользователю варианты действий с выбранным стихотворением
func displayPoemWithOpt(poem *finder.Poem) bool {
	DisplayMessage(fmt.Sprintf("Вы выбрали: %s - %s\nТекст:\n%s\n", poem.Title, poem.Author, poem.Text))

	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("c - Скопировать текст в буфер обмена")
		fmt.Println("Escape - Вернуться к списку стихотворений")

		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return false
		case keyboard.KeyCtrlC:
			copyToClipboard(poem.Text)
			return true
		default:
			continue

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
