package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"pet/finder"
	"pet/terminal"
)

func main() {
	for {
		input, err := terminal.ReadInput()
		// Проверка нажатия Esc (когда возвращается ошибка)
		if err != nil {
			fmt.Println("bye")
			break
		}
		// Проверка на пустую строку (введен только Enter без текста)
		if input == "" {
			if err = termbox.Init(); err != nil {
				panic(err)
			}
			terminal.TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, "Пустой ввод, попробуйте снова. Enter")
			err = termbox.Flush()
			if err != nil {
				return
			}
			termbox.PollEvent()
			termbox.Close()
			continue
		}
		// Обработка введенного текста
		finder.FindPoem(input)

		if len(finder.Poems) == 0 {
			if err = termbox.Init(); err != nil {
				panic(err)
			}
			terminal.TbPrint(0, 0, termbox.ColorDefault, termbox.ColorDefault, "Стихи не найдены.. Enter")
			err = termbox.Flush()
			if err != nil {
				return
			}
			termbox.PollEvent()
			termbox.Close()
			continue
		}
		selectedPoem := terminal.SelectPoem(finder.Poems)
		if selectedPoem == nil {
			continue
		}
	}
}
