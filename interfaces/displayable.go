package interfaces

// Displayable интерфейс для работы с отображением данных в терминале, которые могут отражать заголовок и тело
type Displayable interface {
	TitleDisplayable
	BodyDisplayable
}

// TitleDisplayable интерфейс для структур, которые могут отображать только заголовок
type TitleDisplayable interface {
	DisplayTitle() string
}

// BodyDisplayable интерфейс для структур, которые могут отображать только тело
type BodyDisplayable interface {
	DisplayBody() string
}
