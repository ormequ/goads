package ads

import "fmt"

type Ad struct {
	ID    int64
	Title string
	Text  string
}

func (a Ad) String() string {
	return fmt.Sprintf(
		"<Ad id=%d title=`%s` text=`%s`>",
		a.ID, a.Title, a.Text,
	)
}

func New(title string, text string) Ad {
	return Ad{
		Title: title,
		Text:  text,
	}
}
