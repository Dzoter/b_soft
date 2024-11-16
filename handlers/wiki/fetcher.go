package wiki

const title = "Вики"

type Fetcher struct {
	title string
}

// NewWikiFetcher NewPoemFetcher Используется для main, для определения экземпляра класса
func NewWikiFetcher() Fetcher {
	return Fetcher{
		title: title,
	}
}
func (w Fetcher) DisplayTitle() string {
	return w.title
}
