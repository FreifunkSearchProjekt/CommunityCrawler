package rss

import (
	"bytes"
	"encoding/json"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/config"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/utils"
	"github.com/PuerkitoBio/fetchbot"
	"github.com/mmcdole/gofeed"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type RssFeed struct {
	*sync.WaitGroup
	*config.Config
	*gofeed.Feed
	FC  *fetchbot.Context
	URL *url.URL
}

// FindNewLinks searches for links inside the Feed and crawls them
func (u *RssFeed) FindNewLinks() {
	for _, l := range u.Items {
		u.FC.Q.SendStringHead(l.Link)
	}
	u.Done()
}

// SendData sends the Data to the Indexer
func (u *RssFeed) SendData() {
	defer u.Done()

	log.Println("[INFO] Repacking struct")
	transactionData := &utils.Transaction{}
	transactionData.RssFeed = make([]utils.FeedBasic, 1)
	feedBasic := utils.FeedBasic{
		URL:         u.URL.String(),
		Host:        u.URL.Host,
		Path:        u.URL.Path,
		Title:       u.Title,
		Description: u.Description,
	}
	transactionData.RssFeed[0] = feedBasic

	log.Println("[INFO] Sending to indexer")
	for _, i := range u.Config.Indexer {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(transactionData)
		var url string
		if strings.HasSuffix(i, "/") {
			url = i + "connector_api/index/" + u.Config.CommunityID + "/"
		} else {
			url = i + "/connector_api/index/" + u.Config.CommunityID + "/"
		}

		log.Println("[INFO][INDEXER] Start transaction")
		client := &http.Client{}
		req, reqErr := http.NewRequest("POST", url, b)
		if reqErr != nil {
			log.Println("[ERR][INDEXER] Got error sending: ", reqErr)
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Bearer "+u.Config.CommunityAccessToken)
		res, err := client.Do(req)
		if err != nil {
			log.Println("[ERR][INDEXER] Got error sending: ", err)
		}
		if res.StatusCode != 200 {
			log.Println("[ERR][INDEXER] Some Error occured while contacting indexer: ", res.Status)
		}
	}
}
