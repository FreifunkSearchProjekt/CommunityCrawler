package crawler

import (
	"bytes"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func GetRenderedBody(htm string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return "", err
	}
	bodyFound := cascadia.MustCompile("body").MatchFirst(doc)
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err = html.Render(w, bodyFound)
	if err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
