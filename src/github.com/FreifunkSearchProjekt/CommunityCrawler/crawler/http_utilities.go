package crawler

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
)

func GetRenderedBody(doc *html.Node) (string, error) {
	bn, err := getBody(doc)
	if err != nil {
		return "", err
	}
	body := renderNode(bn)
	return body, nil
}

func getBody(doc *html.Node) (b *html.Node, err error) {
	if doc.Type == html.ElementNode && doc.Data == "body" {
		b = doc
	} else {
		for c := doc.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "body" {
				b = c
			}
		}
	}
	if b == nil {
		err = errors.New("Missing <body> in the node tree")
		return
	}
	return
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
