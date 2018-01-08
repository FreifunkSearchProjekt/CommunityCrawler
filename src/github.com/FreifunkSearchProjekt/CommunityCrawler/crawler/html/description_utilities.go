package html

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/grokify/html-strip-tags-go"
	"strings"
)

func GetDescription(htm *goquery.Document) string {
	// Check if we got a article tag
	article := htm.Find("article").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Html()
		return val
	})

	// If we got a article only search in there for p tags else just get all of them and hope for the best.
	var p []string
	if len(article) > 0 {
		p = htm.Find("article").Find("p").Map(func(_ int, s *goquery.Selection) string {
			val, _ := s.Html()
			return val
		})
	} else {
		//Check if we got role=main
		main := htm.Find("*[role=main]").Map(func(_ int, s *goquery.Selection) string {
			val, _ := s.Html()
			return val
		})

		if len(main) > 0 {
			p = htm.Find("*[role=main]").Find("p").Map(func(_ int, s *goquery.Selection) string {
				val, _ := s.Html()
				return val
			})
		} else {
			// Check if we got a footer and if we can exclude it
			footer := htm.Find("footer").Map(func(_ int, s *goquery.Selection) string {
				val, _ := s.Html()
				return val
			})

			if len(footer) > 0 {
				p = htm.Not("footer").Find("p").Map(func(_ int, s *goquery.Selection) string {
					val, _ := s.Html()
					return val
				})
			} else {
				// Check if we got a header and if we can exclude it
				header := htm.Find("header").Map(func(_ int, s *goquery.Selection) string {
					val, _ := s.Html()
					return val
				})
				if len(header) > 0 {
					p = htm.Not("header").Find("p").Map(func(_ int, s *goquery.Selection) string {
						val, _ := s.Html()
						return val
					})
				} else {
					p = htm.Find("p").Map(func(_ int, s *goquery.Selection) string {
						val, _ := s.Html()
						return val
					})
				}
			}
		}
	}

	if len(p) > 0 {
		joined := strings.Join(p[:], " ")
		stripped := strip.StripTags(joined)
		return stripped
	}
	return ""
}
