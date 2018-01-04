package crawler

import (
	"errors"
	"golang.org/x/net/html"
	"strings"
)

func GetTitle(htm string) (string, error) {
	doc, _ := html.Parse(strings.NewReader(htm))
	tn, err := getTitle(doc)
	if err != nil {
		return "", err
	}
	title := tn.Data
	return title, nil
}

func getTitle(doc *html.Node) (b *html.Node, err error) {
	if checkIfTitle(doc) {
		b = doc
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if checkIfTitle(c) {
			b = c
		}
	}
	if b == nil {
		err = errors.New("Missing <title> in the node tree")
		return
	}
	return
}

func checkIfTitle(n *html.Node) (b bool) {
	if n.Type == html.ElementNode && n.Data == "title" {
		b = true
		return
	}
	return
}
