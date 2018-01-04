package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	CommunityID          string   `yaml:"community_id"`
	CommunityAccessToken string   `yaml:"community_access_token"`
	Indexer              []string `yaml:"indexer"`
	Network              []string `yaml:"network"`
	ExternalPages        []string `yaml:"external_pages"`
}

func loadConfig(filepath string) (config *Config, err error) {
	config = &Config{}
	// detect if file exists
	var _, StatErr = os.Stat(filepath)
	// create file if not exists
	if os.IsNotExist(StatErr) {
		var file, CreateErr = os.Create(filepath)
		if CreateErr != nil {
			err = CreateErr
			return
		}
		defer file.Close()
	}
	fileData, FileErr := ioutil.ReadFile(filepath)
	if FileErr != nil {
		err = FileErr
		return
	}
	YamlErr := yaml.Unmarshal(fileData, config)
	if YamlErr != nil {
		err = YamlErr
		return
	}
	return
}
