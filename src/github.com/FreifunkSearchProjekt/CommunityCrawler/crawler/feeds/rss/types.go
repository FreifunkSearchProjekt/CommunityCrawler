package rss

import (
	"bytes"
	"encoding/json"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/config"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/utils"
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
	URL *url.URL
}

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
		res, err := http.Post(url, "application/json; charset=utf-8", b)
		if err != nil {
			log.Println("[ERR][INDEXER] Got error sending: ", err)
		}
		if res.StatusCode != 200 {
			log.Println("[ERR][INDEXER]Some Error occured while contacting indexer: ", res.Status)
		}
	}
}
