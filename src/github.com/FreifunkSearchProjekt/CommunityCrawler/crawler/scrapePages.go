package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/namsral/microdata"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CrawlFoundings struct {
	UrlsData map[int64]URL
}

type URL struct {
	URL       *url.URL
	Microdata *microdata.Microdata
	Page      string
	Body      string
	Title     string
}

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type Extender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
	CurrentCrawlFoundings   *CrawlFoundings
}

func getPage(url string) (body string, err error) {
	var client http.Client
	resp, GetErr := client.Get(url)
	if GetErr != nil {
		err = GetErr
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, ReadErr := ioutil.ReadAll(resp.Body)
		if ReadErr != nil {
			err = ReadErr
			return
		}
		body = string(bodyBytes)
		return
	}
	return
}

// Override Visit for our need.
func (x *Extender) Visited(ctx *gocrawl.URLContext, harvested interface{}) {
	currentURLData := x.CurrentCrawlFoundings.UrlsData[int64(len(x.CurrentCrawlFoundings.UrlsData)+1)]
	currentURLData.URL = ctx.URL()

	pageMicrodata, err := microdata.ParseURL(ctx.URL().String())
	if err != nil {
		fmt.Errorf("%s", err)
	}
	currentURLData.Microdata = pageMicrodata

	page, err := getPage(ctx.URL().String())
	if err != nil {
		fmt.Errorf("%s", err)
	}
	currentURLData.Page = page

	body, err := GetRenderedBody(currentURLData.Page)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	currentURLData.Body = body

	title, err := GetTitle(currentURLData.Page)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	currentURLData.Title = title
}

func Crawl(url string) (dataToIndex *CrawlFoundings) {
	dataToIndex = &CrawlFoundings{
		UrlsData: make(map[int64]URL),
	}

	extender := new(Extender)
	extender.CurrentCrawlFoundings = dataToIndex

	// Set custom options
	opts := gocrawl.NewOptions(extender)

	// should always set your robot name so that it looks for the most
	// specific rules possible in robots.txt.
	opts.RobotUserAgent = "FreifunkSearchProjektCrawler"

	opts.LogFlags = gocrawl.LogAll
	opts.SameHostOnly = true

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(url)

	dataToIndex = extender.CurrentCrawlFoundings
	return
}
