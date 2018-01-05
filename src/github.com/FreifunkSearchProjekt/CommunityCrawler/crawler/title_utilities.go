package crawler

import (
	"bytes"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

var stripTitle = regexp.MustCompile(`(?:<title>)(.*)(?:<\/title>)`)

func GetTitle(htm string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return "", err
	}
	titleFound := cascadia.MustCompile("title").MatchFirst(doc)

	var buf bytes.Buffer
	w := io.Writer(&buf)
	RenderErr := html.Render(w, titleFound)
	if RenderErr != nil {
		return "", RenderErr
	}
	titleNode := buf.String()
	var title string
	if len(titleNode) > 0 {
		title = stripTitle.FindAllStringSubmatch(titleNode, -1)[0][1]
	}

	return title, nil
}
