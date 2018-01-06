package html

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/grokify/html-strip-tags-go"
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

	stripped := strip.StripTags(title)
	return stripped, nil
}
