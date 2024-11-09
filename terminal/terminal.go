package terminal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"pet/interfaces"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/mattn/go-tty"
	_ "pet/interfaces"
)

const (
	ActionNext = "next"
	ActionPrev = "prev"
	ActionExit = "exit"
)

// ReadInput считывает ввод пользователя (текст)
func ReadInput() string {
	//TODO реализовать многопоточку что бы отловить нажатие escape и выйти из метода
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите текст")

	fmt.Print("> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении ввода:", err)
		return ""
	}

	input = strings.TrimSpace(input)
	return input
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
func DisplayItemTitle(item interfaces.Displayable) {
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
func SelectItemsWithoutPaging(items []interfaces.Displayable) (interfaces.Displayable, string) {
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
