package terminal

import (
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/nsf/termbox-go"
	"pet/finder"
	"strings"
)

// ReadInput считывает ввод пользователя
func ReadInput() (string, error) {

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Очистить экран и вывести приглашение
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		return "", err
	}
	TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, "Введите текст стиха для поиска (Esc - выйти):")
	err = termbox.Flush()
	if err != nil {
		return "", err
	}

	var input strings.Builder
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return "", errors.New("user exited") // Если нажата клавиша Esc, возвращаем пустую строку для выхода из программы
			case termbox.KeyEnter:
				termbox.Close()
				return strings.TrimSpace(input.String()), nil // Возвращаем введенный текст без пробелов
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if input.Len() > 0 {
					// Конвертируем строку в срез рун для корректной работы с юникод символами
					runes := []rune(input.String())
					// Удаляем последний символ
					input.Reset()
					input.WriteString(string(runes[:len(runes)-1]))
				}
			default:
				if ev.Ch != 0 {
					input.WriteRune(ev.Ch) // Добавляем введенный символ к строке
				}
			}

			// Обновляем экран с текущим вводом
			err = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			if err != nil {
				return "", err
			}
			TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, "Введите текст стиха для поиска (Esc - выйти):")
			TbPrint(0, 1, termbox.ColorDefault, termbox.ColorDefault, input.String()+"\n")
			err = termbox.Flush()
			if err != nil {
				return "", err
			}
		default:
			continue
		}
	}
}

// displayPoems отображает список найденных стихов
func displayPoems(poems []finder.Poem, currentIndex int) {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		return
	}
	TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, "Найденные стихи:")
	for i, poem := range poems {
		line := fmt.Sprintf("%-2d %s - %s", i, poem.Title, poem.Author)
		color := termbox.ColorDefault
		if i == currentIndex {
			color = termbox.ColorCyan
		}
		TbPrint(0, i+1, color, termbox.ColorDefault, line)
	}
	err = termbox.Flush()
	if err != nil {
		return
	}
}

// TbPrint formats and prints a string to the termbox buffer
func TbPrint(x, y int, fg, bg termbox.Attribute, msg string) int {
	for _, c := range msg {
		if c == '\n' {
			y++
			x = 0
		} else if c == '\r' {
			continue // игнорируем '\r', он часто используется вместе с '\n'
		} else {
			termbox.SetCell(x, y, c, fg, bg)
			x++
		}
	}
	return y
}

// SelectPoem предлагает пользователю выбрать стихотворение из списка
func SelectPoem(poems []finder.Poem) *finder.Poem {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	currentIndex := 0

	for {
		displayPoems(poems, currentIndex)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if currentIndex > 0 {
					currentIndex--
				}
			case termbox.KeyArrowDown:
				if currentIndex < len(poems)-1 {
					currentIndex++
				}
			case termbox.KeyEnter:
				displayPoemOptions(poems[currentIndex])
				continue
			case termbox.KeyEsc:
				return nil
			default:
				continue
			}
		case termbox.EventError:
			return nil
		default:
			continue
		}
	}
}

func displayPoemOptions(poem finder.Poem) *finder.Poem {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		return nil
	}

	// Выводим заголовок и текст стихотворения
	lastLine := TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, fmt.Sprintf("Вы выбрали: %s - %s\nТекст:\n%s", poem.Title, poem.Author, poem.Text))

	// Рассчитываем позицию для следующего блока текста
	lastLine++

	TbPrint(0, lastLine, termbox.ColorDefault, termbox.ColorDefault, "\nВыберите действие:")
	lastLine++
	TbPrint(0, lastLine, termbox.ColorDefault, termbox.ColorDefault, "С - Скопировать текст в буфер обмена")
	lastLine++
	TbPrint(0, lastLine, termbox.ColorDefault, termbox.ColorDefault, "Esc - Вернуться к списку стихотворений")
	err = termbox.Flush()
	if err != nil {
		return nil
	}

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Key {
			case termbox.KeyEsc:
				return nil
			default:
				switch ev.Ch {
				case 'c', 'C', 'с', 'С':
					copyToClipboard(poem.Text)
					return nil
				default:
					err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
					if err != nil {
						return nil
					}
					TbPrint(0, lastLine+1, termbox.ColorRed, termbox.ColorDefault, "Неверный выбор, попробуйте снова.")
					err = termbox.Flush()
					if err != nil {
						return nil
					}
				}
			}
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
