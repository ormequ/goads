package links

import (
	"fmt"
)

type Link struct {
	ID       int64
	URL      string `validate:"min:1;max:2048"`
	Alias    string `validate:"min:1"`
	AuthorID int64
	Ads      []int64
}

func (l Link) String() string {
	return fmt.Sprintf(
		"<Ad id=%d authorID=%d url=`%s` alias=`%v` ads=%v>",
		l.ID,
		l.AuthorID,
		l.URL,
		l.Alias,
		l.Ads,
	)
}

func New(url string, alias string, authorID int64, ads []int64) Link {
	return Link{
		URL:      url,
		Alias:    alias,
		AuthorID: authorID,
		Ads:      ads,
	}
}
