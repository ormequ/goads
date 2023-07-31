package redirects

import (
	"fmt"
	"goads/internal/urlshortener/entities/ads"
	"goads/internal/urlshortener/entities/links"
)

type Redirect struct {
	Link links.Link
	Ad   ads.Ad
}

func (r Redirect) String() string {
	return fmt.Sprintf("<Redirect link=`%s` ad=`%s`>", r.Link.String(), r.Ad.String())
}

func New(link links.Link, ad ads.Ad) Redirect {
	return Redirect{
		Link: link,
		Ad:   ad,
	}
}
