package interfaces

import "pet/structs"

type Actionable interface {
	GetActions() []structs.KeyBinding // Возвращает список доступных действий
}
