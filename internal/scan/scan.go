package scan

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/babilon15/trfeed/internal/addtorrent"
	"github.com/babilon15/trfeed/internal/config"
	"github.com/babilon15/trfeed/pkg/feed"
	"github.com/babilon15/trfeed/pkg/utils"
)

const (
	configFile   = "config.yaml"
	remnantsFile = "remnants.yaml"
)

var (
	hits Hits
	conf config.Config
)

func init() {
	if err := utils.GetYAMLFromFile(remnantsFile, &hits); err != nil {
		fmt.Println(err)
	}

	if err := utils.GetYAMLFromFile(configFile, &conf); err != nil {
		fmt.Println(err)
	}

	if conf.TargetDir == "" {
		fmt.Println("[Warn] target_dir is missing from", configFile)

		dd := addtorrent.GetDownDirWithRemote(conf.Host, conf.Auth)
		if dd == "" {
			fmt.Println("[Error] The client's default download directory is unknown. The host may be unavailable.")
		} else {
			conf.TargetDir = dd
			fmt.Println("[Info] New 'target_dir' directory from the client:", dd)
		}
	}
}

func GetFeed(url string, target any) error {
	resp, respErr := http.Get(url)
	if respErr != nil {
		return respErr
	}
	defer resp.Body.Close()

	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return bodyErr
	}

	return xml.Unmarshal(body, target)
}

func Scan() {
	feedsLen := len(conf.Feeds)

	if feedsLen == 0 {
		fmt.Println("[Warn] No feeds could be found in the config file:", configFile)
	}

	for i := 0; i < feedsLen; i++ {
		if !conf.Feeds[i].GetAll && len(conf.Feeds[i].Filters) == 0 && len(conf.Feeds[i].FiltersViaLabels) == 0 {
			continue
		}

		if !utils.IsValidURL(conf.Feeds[i].Url) {
			fmt.Println("Invalid URL:", conf.Feeds[i].Url)
			continue
		}

		var currentFeed feed.Feed

		if err := GetFeed(conf.Feeds[i].Url, &currentFeed); err != nil {
			fmt.Println(err)
			continue
		}

	outer:
		for _, v := range currentFeed.Channel.Item {
			// (1) own filters
			for _, f := range conf.Feeds[i].Filters {
				if f.Check(v.Title) {
					// HIT!
					fmt.Println("[Hit]", v.Title)
					//
					continue outer
				}
			}

			// (2) filters via labels
			for _, l := range conf.Feeds[i].FiltersViaLabels {
				currentFilter := conf.GetFilterByLabel(l)
				if config.IsFilterEmpty(currentFilter) {
					continue
				}

				if currentFilter.Check(v.Title) {
					// HIT!
					fmt.Println("[Hit]", v.Title)
					//
					continue outer
				}
			}

			// (3) get all
			if conf.Feeds[i].GetAll {
				//
			}
		}
	}

	if err := utils.PutYAMLToFile(remnantsFile, &hits); err != nil {
		fmt.Println(err)
	}
}

func AddHits() {
	for i := 0; i < len(hits); i++ {
		fmt.Println(hits[i])
	}
}
