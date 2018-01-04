package common

import (
	"bytes"
	"encoding/json"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/crawler"
	"log"
	"net/http"
	"strings"
)

func Setup(configPath string) (*Config, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Begin(config *Config) {
	//First crawl external
	log.Println("Crawl External Pages")
	for _, e := range config.ExternalPages {
		work(e, config)
	}

	// After it Crawl the network
	log.Println("Crawl Networks")
	for _, i := range config.Network {
		//TODO Handle error
		hosts, _ := crawler.Hosts(i)
		for _, h := range hosts {
			for _, p := range crawler.ScanServer(h) {
				work(h+":"+string(p), config)
			}
		}
	}
}

func work(url string, config *Config) {
	results := crawler.Crawl(url)

	transactionData := transaction{}
	transactionData.BasicWebpages = make([]WebpageBasic, len(results.UrlsData))
	for i, u := range results.UrlsData {
		page := WebpageBasic{
			URL:   u.URL.String(),
			Path:  u.URL.Path,
			Title: u.Title,
			Body:  u.Body,
		}
		transactionData.BasicWebpages[i] = page
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(transactionData)
	log.Println(b.String())
	for _, i := range config.Indexer {
		var url string
		if strings.HasSuffix(i, "/") {
			url = i + "connector_api/index/" + config.CommunityID + "/"
		} else {
			url = i + "/connector_api/index/" + config.CommunityID + "/"
		}

		log.Println(url)
		res, err := http.Post(url, "application/json; charset=utf-8", b)
		if res.StatusCode != 200 {
			log.Println("Some Error occured while contacting indexer: ", res.Status)
		}
		if err != nil {
			log.Println("Got error sending: ", err)
		}
	}
}
