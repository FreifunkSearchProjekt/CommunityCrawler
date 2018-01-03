package crawler

import (
	"bytes"
	"golang.org/x/net/html"
	"errors"
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

func getBody(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}