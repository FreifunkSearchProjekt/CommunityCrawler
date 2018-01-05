package crawler

import (
	"bytes"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"io"
	"strings"
)

func GetRenderedBody(htm string) (string, error) {
	var body string
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return "", err
	}
	bodyFound := cascadia.MustCompile("body").MatchFirst(doc)

	if bodyFound != nil {
		var buf bytes.Buffer
		w := io.Writer(&buf)
		RenderErr := html.Render(w, bodyFound)
		if RenderErr != nil {
			return "", RenderErr
		}
		body = buf.String()
	}

	return body, nil
}
