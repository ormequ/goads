package ads

import (
	"fmt"
	"time"
)

type Ad struct {
	ID         int64
	AuthorID   int64
	Published  bool
	Title      string `validate:"min:1; max:99"`
	Text       string `validate:"min:1"`
	CreateDate time.Time
	UpdateDate time.Time
}

func (a Ad) String() string {
	return fmt.Sprintf(
		"<Ad id=%d authorID=%d created=%s published=%v title=`%s` text=`%s`>",
		a.ID,
		a.AuthorID,
		a.CreateDate,
		a.Published,
		a.Title,
		a.Text,
	)
}

func (a Ad) GetCreateDate() time.Time {
	return a.CreateDate.Truncate(24 * time.Hour)
}

func (a Ad) GetUpdateDate() time.Time {
	return a.UpdateDate.Truncate(24 * time.Hour)
}

func New(title string, text string, authorID int64) Ad {
	return Ad{
		AuthorID:   authorID,
		Title:      title,
		Text:       text,
		CreateDate: time.Now().UTC(),
		UpdateDate: time.Now().UTC(),
		Published:  true,
	}
}
