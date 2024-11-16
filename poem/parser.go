package poem

const title = "Стихи"

type Parser struct {
	title string
}

// NewPoemParser Используется для main, для определения экземпляра класса
func NewPoemParser() Parser {
	return Parser{
		title: title,
	}
}
func (p Parser) DisplayTitle() string {
	return p.title
}
func (p Parser) DisplayBody() string {
	return ""
}
