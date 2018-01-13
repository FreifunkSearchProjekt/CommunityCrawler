package html

import (
	"bytes"
	"encoding/json"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/config"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/utils"
	"github.com/namsral/microdata"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type URL struct {
	*sync.WaitGroup
	*config.Config
	URL         *url.URL
	Microdata   *microdata.Microdata
	Body        string
	Title       string
	Description string
}

// SendData sends the Data to the Indexer
func (u *URL) SendData() {
	defer u.Done()

	log.Println("[INFO] Repacking struct")
	transactionData := utils.Transaction{}
	transactionData.BasicWebpages = make([]utils.WebpageBasic, 1)
	webpageBasic := utils.WebpageBasic{
		URL:         u.URL.String(),
		Host:        u.URL.Host,
		Path:        u.URL.Path,
		Title:       u.Title,
		Body:        u.Body,
		Description: u.Description,
	}
	transactionData.BasicWebpages[0] = webpageBasic

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
