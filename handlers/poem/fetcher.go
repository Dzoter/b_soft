package poem

const title = "Стихи"

type Fetcher struct {
	title string
}

// NewPoemFetcher Используется для main, для определения экземпляра класса
func NewPoemFetcher() Fetcher {
	return Fetcher{
		title: title,
	}
}
func (p Fetcher) DisplayTitle() string {
	return p.title
}
