package html

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/grokify/html-strip-tags-go"
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

	stripped := strip.StripTags(body)
	return stripped, nil
}
