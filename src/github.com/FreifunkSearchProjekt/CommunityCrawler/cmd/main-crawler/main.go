package main

import (
	"flag"
	"github.com/FreifunkSearchProjekt/CommunityCrawler/common"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "crawler.yaml", "Define the config location")

	flag.Parse()
	config, err := common.Setup(configPath)
	if err != nil {
		log.Panic(err)
	}

	common.Begin(config)
}
