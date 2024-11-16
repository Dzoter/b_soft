package terminal

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/mattn/go-tty"
	"log"
	"os"
	"os/exec"
	"pet/interfaces"
	_ "pet/interfaces"
	"runtime"
	"strings"
)

const (
	ActionNext = "next"
	ActionPrev = "prev"
	ActionExit = "exit"
)

// ReadInput считывает ввод пользователя (текст)
func ReadInput() string {
	// Создаем переменную для накопления введённых символов
	var inputBuilder strings.Builder

	// Сообщение для пользователя
	DisplayMessage("Введите текст (нажмите Enter для подтверждения, Esc для выхода)")

	// Открываем клавиатуру для отслеживания ввода
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	fmt.Print("> ")

	// Бесконечный цикл для обработки нажатий клавиш
	for {
		// Получаем нажатую клавишу
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			fmt.Println("\nВыход из метода по сигналу Escape")
			return ActionExit // Возвращаем специальное значение для выхода

		case keyboard.KeyEnter:
			// Возвращаем накопленный текст при нажатии Enter
			fmt.Println() // Переход на новую строку после ввода
			return inputBuilder.String()

		case keyboard.KeyBackspace2:
			// Удаляем последний символ при нажатии Backspace
			currentInput := inputBuilder.String()
			runes := []rune(currentInput) // Преобразуем строку в срез рун

			if len(runes) > 0 {
				// Убираем последний символ из среза рун
				runes = runes[:len(runes)-1]
				// Перезаписываем строку без последнего символа
				inputBuilder.Reset()
				inputBuilder.WriteString(string(runes))
				fmt.Print("\b \b") // Убираем символ с экрана
			}

		default:
			// Добавляем введённый символ к накопленному тексту и выводим его
			if char != 0 { // Проверяем, что это символ, а не специальная клавиша
				inputBuilder.WriteRune(char)
				fmt.Print(string(char)) // Выводим символ на экран
			}
		}
	}
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Print("\033[H\033[2J")
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// DisplayItemTitle DisplayItem отображает элемент
func DisplayItemTitle(item interfaces.TitleDisplayable) {
	clearScreen()
	DisplayMessage(item.DisplayTitle())
	moveCursorUp()
}

// DisplayItemBody DisplayItemTitle DisplayItem отображает элемент
func DisplayItemBody(item interfaces.Displayable) {
	clearScreen()
	DisplayMessage(item.DisplayBody())
	moveCursorUp()
}

func DisplayMessage(message string) {
	clearScreen()
	fmt.Printf("%s\n", message)
}

// SelectItemsWithPaging позволяет пользователю выбирать элементы с использованием клавиш вверх/вниз
func SelectItemsWithPaging(items []interfaces.Displayable) (interfaces.Displayable, string) {
	// Инициализация клавиатуры
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	choice := -1 // Начинаем выбор с нулевого элемента

	DisplayMessage("Объекты загружены, нажимайте вверх/вниз для выбора")

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return nil, "exit"

		case keyboard.KeyArrowUp:
			if choice > 0 {
				choice-- // Перемещаем выбор вверх
				DisplayItemTitle(items[choice])
			} else {
				return nil, "prev"
			}

		case keyboard.KeyArrowDown:
			if choice < len(items)-1 {
				choice++ // Перемещаем выбор вниз
				DisplayItemTitle(items[choice])
			} else {
				return nil, "next"
			}

		case keyboard.KeyEnter:
			return items[choice], ""

		default:
			continue
		}
	}
}

// SelectItemsWithoutPaging предлагает пользователю выбрать стрелочками элемент из списка
func SelectItemsWithoutPaging(items []interfaces.TitleDisplayable) (interfaces.TitleDisplayable, string) {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	choice := -1 // начальное значение choice, отрицательное что-бы в начале всегда отображался 0 элемент

	DisplayMessage("Объекты загружены, нажимайте вверх/вниз для выбора")

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyEsc:
			return nil, ActionExit
		case keyboard.KeyArrowUp:
			if choice > 0 {
				choice--
				DisplayItemTitle(items[choice])
			}
		case keyboard.KeyArrowDown:
			if choice < len(items)-1 {
				choice++
				DisplayItemTitle(items[choice])
			}
		case keyboard.KeyEnter:
			return items[choice], ""
		default:
			continue
		}
	}
}

func moveCursorUp() {
	t, err := tty.Open()
	if err != nil {
		fmt.Println("Ошибка при открытии TTY:", err)
		return
	}
	defer t.Close()

	// Отправляем ANSI escape-код для перемещения курсора вверх
	t.Output().WriteString("\033[A")

	// Принудительное обновление
	t.Output().Sync()
}

func DisplayActions(item interfaces.Actionable) bool {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	bindings := item.GetActions()

	for {
		fmt.Println("\nВыберите действие:")
		for _, binding := range bindings {
			fmt.Printf("%s - %s\n", string(binding.Char), binding.Description)
		}
		fmt.Println("Escape - Вернуться к списку")

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		// Проверка выхода
		if key == keyboard.KeyEsc {
			return false
		}

		// Поиск и выполнение действия по нажатой клавише
		for _, binding := range bindings {
			if key == binding.Key || char == binding.Char {
				return binding.Action()
			}
		}
	}
}
