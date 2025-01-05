package interfaces

import "b_soft/structs"

type Actionable interface {
	GetActions() []structs.KeyBinding // Возвращает список доступных действий
}
