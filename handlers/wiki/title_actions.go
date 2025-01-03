package wiki

import "fmt"

func (t Title) DisplayTitle() string {
	return fmt.Sprintf("%s \n", t.Title)
}

func (t Title) DisplayBody() string {
	return fmt.Sprintf("%s \n", "")
}
