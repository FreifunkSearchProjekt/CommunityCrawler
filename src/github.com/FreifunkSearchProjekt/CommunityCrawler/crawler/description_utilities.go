package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/grokify/html-strip-tags-go"
)

func GetDescription(htm *goquery.Document) string {
	p := htm.Find("p:first-of-type").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Html()
		return val
	})

	if len(p) > 0 {
		stripped := strip.StripTags(p[0])
		return stripped
	}
	return ""
}
