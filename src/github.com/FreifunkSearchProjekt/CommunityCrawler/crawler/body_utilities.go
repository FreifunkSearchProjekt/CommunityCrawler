package crawler

import (
	"github.com/PuerkitoBio/goquery"
)

func GetRenderedBody(htm *goquery.Document) (string, error) {
	var body string
	var err error

	htm.Find("body").Each(func(i int, s *goquery.Selection) {
		body, err = s.Html()
	})

	if err != nil {
		return "", err
	}

	return body, nil
}
