package common

import (
	"github.com/FreifunkSearchProjekt/CommunityCrawler/config"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/crawler"
	"log"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/crawler/scan"
)

func Setup(configPath string) (*config.Config, error) {
	configData, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return configData, nil
}

func Begin(config *config.Config) {
	//First crawl external
	log.Println("Crawl External Pages")
	for _, e := range config.ExternalPages {
		work(e, config)
	}

	// After it Crawl the network
	log.Println("Crawl Networks")
	for _, i := range config.Network {
		//TODO Handle error
		hosts, _ := scan.Hosts(i)
		for ip, ports := range scan.Scan(hosts) {
			for _, p := range ports {
				work(ip+":"+string(p), config)
			}
		}
	}
}

func work(url string, config *config.Config) {
	crawler.Crawl(url, config)
}
