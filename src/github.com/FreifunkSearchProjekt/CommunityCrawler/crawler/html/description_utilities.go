package html

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/grokify/html-strip-tags-go"
	"strings"
)

func GetDescription(htm *goquery.Document) string {
	p := htm.Find("p:first-of-type").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Html()
		return val
	})

	if len(p) > 0 {
		joined := strings.Join(p[:], " ")
		stripped := strip.StripTags(joined)
		return stripped
	}
	return ""
}
