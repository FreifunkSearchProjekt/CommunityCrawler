package crawler

import (
	"bytes"
	"encoding/json"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/config"
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

	purellFlags purell.NormalizationFlags
)

type URL struct {
	URL         *url.URL
	Microdata   *microdata.Microdata
	Body        string
	Title       string
	Description string
}

func Crawl(urlS string, config *config.Config) {
	// Don't force HTTP
	purellFlags = purell.FlagDecodeDWORDHost | purell.FlagDecodeOctalHost | purell.FlagDecodeHexHost | purell.FlagRemoveUnnecessaryHostDots | purell.FlagRemoveEmptyPortSeparator | purell.FlagsUsuallySafeGreedy | purell.FlagRemoveDirectoryIndex | purell.FlagRemoveFragment | purell.FlagRemoveDuplicateSlashes

	// Parse the provided url
	normalized, err := purell.NormalizeURLString(urlS, purellFlags)
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
	// TODO refactor
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

			currentURLData := &URL{}
			currentURLData.URL = ctx.Cmd.URL()

			pageMicrodata, err := microdata.ParseURL(ctx.Cmd.URL().String())
			if err != nil {
				log.Printf("[ERR] %s", err)
			}
			currentURLData.Microdata = pageMicrodata

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
			if err != nil {
				log.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
				return
			}

			body, err := GetRenderedBody(doc)
			if err != nil {
				log.Printf("[ERR] %s", err)
			}
			if len(body) > 0 {
				currentURLData.Body = body
			} else {
				currentURLData.Body = page
			}

			title, err := GetTitle(doc)
			if err != nil {
				log.Printf("[ERR] %s", err)
			}
			if len(title) > 0 {
				currentURLData.Title = title
			} else {
				currentURLData.Title = ctx.Cmd.URL().String()
			}

			description := GetDescription(doc)
			currentURLData.Description = description

			//Send Data
			transactionData := transaction{}
			transactionData.BasicWebpages = make([]WebpageBasic, 1)
			webpageBasic := WebpageBasic{
				URL:         currentURLData.URL.String(),
				Path:        currentURLData.URL.Path,
				Title:       currentURLData.Title,
				Body:        currentURLData.Body,
				Description: currentURLData.Description,
			}
			transactionData.BasicWebpages[0] = webpageBasic

			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(transactionData)
			for _, i := range config.Indexer {
				var url string
				if strings.HasSuffix(i, "/") {
					url = i + "connector_api/index/" + config.CommunityID + "/"
				} else {
					url = i + "/connector_api/index/" + config.CommunityID + "/"
				}

				_, err := http.Post(url, "application/json; charset=utf-8", b)
				/*		if res.StatusCode != 200 {
						log.Println("Some Error occured while contacting indexer: ", res.Status)
					}*/
				if err != nil {
					log.Println("Got error sending: ", err)
				}
			}

			// Process the body to find the links
			// Enqueue all links as HEAD requests
			enqueueLinks(ctx, doc, u)
			return
		}))

	// Handle HEAD requests for html responses coming from the source host - we don't want
	// to crawl links from other hosts.
	mux.Response().Method("HEAD").Host(u.Host).ContentType("text/html").Handler(fetchbot.HandlerFunc(
		func(ctx *fetchbot.Context, res *http.Response, err error) {
			if _, err := ctx.Q.SendStringGet(ctx.Cmd.URL().String()); err != nil {
				log.Printf("[ERR] %s %s - %s\n", ctx.Cmd.Method(), ctx.Cmd.URL(), err)
			}
			return
		}))

	// Create the Fetcher, handle the logging first, then dispatch to the Muxer
	h := logHandler(mux)
	f := fetchbot.New(h)
	f.AutoClose = true
	f.UserAgent = "FreifunkSearchProjektCrawler"

	// Start processing
	q := f.Start()

	// Enqueue the seed, which is the first entry in the dup map
	dup[normalized] = true
	_, err = q.SendStringGet(normalized)
	if err != nil {
		log.Printf("[ERR] GET %s - %s\n", normalized, err)
	}
	q.Block()
	return
}

// logHandler prints the fetch information and dispatches the call to the wrapped Handler.
func logHandler(wrapped fetchbot.Handler) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		if err == nil {
			log.Printf("[%d] %s %s - %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL(), res.Header.Get("Content-Type"))
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
	normalized, err := purell.NormalizeURLString(resolvedURL.String(), purellFlags)
	if err != nil {
		log.Printf("error: normalize URL %s - %s\n", resolvedURL.String(), err)
	}
	return normalized
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
			normalized, err := purell.NormalizeURLString(s, purellFlags)
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
