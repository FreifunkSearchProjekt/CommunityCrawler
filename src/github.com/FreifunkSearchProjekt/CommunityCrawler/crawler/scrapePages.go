package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
	"github.com/namsral/microdata"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

var (
	// Protect access to dup
	mu sync.Mutex

	// Duplicates table
	dup = map[string]bool{}
)

type URL struct {
	URL       *url.URL
	Microdata *microdata.Microdata
	Body      string
	Title     string
}

func Crawl(urlS string) (dataToIndex map[int64]*URL) {

	var UrlsData = make(map[int64]*URL)

	// Parse the provided url
	normalized, err := purell.NormalizeURLString(urlS, purell.FlagsAllGreedy)
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(normalized)
	if err != nil {
		log.Fatal(err)
	}

	// Create the muxer
	mux := fetchbot.NewMux()

	// Handle all errors the same
	mux.HandleErrors(fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		log.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
	}))

	// Handle GET requests for html responses, to parse the body and enqueue all links as HEAD
	// requests.
	mux.Response().Method("GET").ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			//Process current URL
			var page string
			defer res.Body.Close()

			pageBytes, ReadErr := ioutil.ReadAll(res.Body)
			if ReadErr != nil {
				err = ReadErr
				return
			}
			page = string(pageBytes)

			log.Println(len(UrlsData))
			currentURLData := &URL{}
			currentURLData.URL = ctx.Cmd.URL()

			pageMicrodata, err := microdata.ParseURL(ctx.Cmd.URL().String())
			if err != nil {
				fmt.Errorf("%s", err)
			}
			currentURLData.Microdata = pageMicrodata

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
			if err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				return
			}

			body, err := GetRenderedBody(doc)
			if err != nil {
				fmt.Errorf("%s", err)
			}
			if len(body) > 0 {
				currentURLData.Title = body
			} else {
				currentURLData.Title = page
			}

			title, err := GetTitle(doc)
			if err != nil {
				fmt.Errorf("%s", err)
			}
			if len(title) > 0 {
				currentURLData.Title = title
			} else {
				currentURLData.Title = ctx.Cmd.URL().String()
			}
			UrlsData[int64(len(UrlsData))] = currentURLData

			// Process the body to find the links
			// Enqueue all links as HEAD requests
			enqueueLinks(ctx, doc, u)
		}))

	// Handle HEAD requests for html responses coming from the source host - we don't want
	// to crawl links from other hosts.
	mux.Response().Method("HEAD").Host(u.Host).ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			if _, err := ctx.Q.SendStringGet(ctx.Cmd.URL().String()); err != nil {
				fmt.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
			}
		}))

	// Create the Fetcher, handle the logging first, then dispatch to the Muxer
	h := logHandler(mux)
	f := fetchbot.New(h)

	// Start processing
	q := f.Start()

	// Enqueue the seed, which is the first entry in the dup map
	dup[urlS] = true
	_, err = q.SendStringGet(urlS)
	if err != nil {
		fmt.Printf("[ERR] GET %s - %s\n", urlS, err)
	}
	q.Block()

	dataToIndex = UrlsData
	return
}

// logHandler prints the fetch information and dispatches the call to the wrapped Handler.
func logHandler(wrapped fetchbot.Handler) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		if err == nil {
			if ctx.Cmd.Method() != "HEAD" {
				fmt.Printf("[%d] %s %s - %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL(), res.Header.Get("Content-Type"))
			}
		}
		wrapped.Handle(ctx, res, err)
	})
}

func handleBaseTag(root *url.URL, baseHref string, aHref string) string {
	resolvedBase, err := root.Parse(baseHref)
	if err != nil {
		return ""
	}

	parsedURL, err := url.Parse(aHref)
	if err != nil {
		return ""
	}
	// If a[href] starts with a /, it overrides the base[href]
	if parsedURL.Host == "" && !strings.HasPrefix(aHref, "/") {
		aHref = path.Join(resolvedBase.Path, aHref)
	}

	resolvedURL, err := resolvedBase.Parse(aHref)
	if err != nil {
		return ""
	}
	return resolvedURL.String()
}

func enqueueLinks(ctx *fetchbot.Context, doc *goquery.Document, originalURL *url.URL) {
	mu.Lock()
	baseURL, _ := doc.Find("base[href]").Attr("href")
	urls := doc.Find("a[href]").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Attr("href")
		if baseURL != "" {
			val = handleBaseTag(doc.Url, baseURL, val)
		}
		return val
	})

	for _, s := range urls {
		if len(s) > 0 && !strings.HasPrefix(s, "#") {
			// Resolve address
			normalized, err := purell.NormalizeURLString(s, purell.FlagsAllGreedy)
			if err != nil {
				log.Printf("error: normalize URL %s - %s\n", s, err)
			}
			u, err := ctx.Cmd.URL().Parse(normalized)
			if err != nil {
				log.Printf("error: resolve URL %s - %s\n", s, err)
				return
			}
			// If prevents sending unnecessary Head requests
			// Ignore URLs that have a #
			// Ignore URLs that have ?
			// Ignore URLs with different scheme than https or http
			// Ignore if host of href isn't the same as the original host
			if u.Fragment == "" && u.RawQuery == "" && (u.Scheme == "https" || u.Scheme == "http") && u.Host == originalURL.Host {
				if !dup["http://"+u.Host+u.Path] && !dup["https://"+u.Host+u.Path] {
					if _, err := ctx.Q.SendStringHead(u.String()); err != nil {
						log.Printf("error: enqueue head %s - %s\n", u, err)
					} else {
						dup[u.String()] = true
					}
				}
			} else {
				// If prevents sending unnecessary Head requests
				// Ignore if already duplicated
				// Ignore URLs with different scheme than https or http
				// Ignore if host of href isn't the same as the original host
				if !dup["http://"+u.Host+u.Path] && !dup["https://"+u.Host+u.Path] && (u.Scheme == "https" || u.Scheme == "http") && u.Host == originalURL.Host {
					if _, err := ctx.Q.SendStringHead(u.Scheme + "://" + u.Host + u.Path); err != nil {
						log.Printf("error: enqueue head %s - %s\n", u, err)
					} else {
						dup[u.Scheme+"://"+u.Host+u.Path] = true
					}
				}
				// Index the without fragment version if not done before
				dup[u.String()] = true
			}
		}
	}
	mu.Unlock()
	return
}
