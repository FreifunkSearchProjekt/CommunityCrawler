package crawler

import (
	"github.com/PuerkitoBio/goquery"
)

func GetTitle(htm *goquery.Document) (string, error) {
	var title string
	var err error

	htm.Find("title").Each(func(i int, s *goquery.Selection) {
		title, err = s.Html()
	})
	if err != nil {
		return "", err
	}

	return title, nil
}
