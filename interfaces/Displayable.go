package interfaces

// Displayable интерфейс для работы с отображением данных в терминале
type Displayable interface {
	DisplayTitle() string
	DisplayBody() string
}
