package scan

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
		fmt.Println("target_dir is missing from", configFile)

		dd := addtorrent.GetDownDirWithRemote(conf.Host, conf.Auth)
		if dd == "" {
			fmt.Println("the client's default download directory is unknown; the host may be unavailable")
		} else {
			conf.TargetDir = dd
			fmt.Println("new 'target_dir' directory from the client:", dd)
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

func handleTargetDirs(dirs ...string) string {
	for _, d := range dirs {
		if d == "" {
			continue
		}

		err := os.MkdirAll(d, 0o777)
		if err != nil {
			fmt.Println(err)
		} else {
			return d
		}
	}
	return ""
}

func Scan() {
	feedsLen := len(conf.Feeds)

	if feedsLen == 0 {
		fmt.Println("no feeds could be found in the config file:", configFile)
	}

	for i := 0; i < feedsLen; i++ {
		if !conf.Feeds[i].GetAll && len(conf.Feeds[i].Filters) == 0 && len(conf.Feeds[i].FiltersViaLabels) == 0 {
			continue
		}

		if !utils.IsValidURL(conf.Feeds[i].Url) {
			fmt.Println("invalid URL:", conf.Feeds[i].Url)
			continue
		}

		var currentFeed feed.Feed

		if err := GetFeed(conf.Feeds[i].Url, &currentFeed); err != nil {
			fmt.Println(err)
			continue
		}

		var firstUniqueNum uint32

	outer:
		for j, v := range currentFeed.Channel.Item {
			currentUniqueNum := v.GetUniqueNum()
			if j == 0 {
				firstUniqueNum = currentUniqueNum
			}

			if conf.Feeds[i].LastUniqueNum != 0 {
				if conf.Feeds[i].LastUniqueNum == currentUniqueNum {
					conf.Feeds[i].LastUniqueNum = firstUniqueNum
					break
				}
			} else {
				conf.Feeds[i].LastUniqueNum = firstUniqueNum
			}

			if hits.IndexByUniqueNum(currentUniqueNum) != -1 {
				continue outer
			}

			// (1) own filters
			for _, f := range conf.Feeds[i].Filters {
				if f.Check(v.Title) {
					// HIT!
					fmt.Println("hit:", v.Title)

					hits = append(hits, Hit{
						Labels:    utils.FilterEmptyStrings([]string{"trfeed", f.Label}),
						Title:     v.Title,
						Resource:  v.Link,
						TargetDir: handleTargetDirs(f.TargetDir, conf.Feeds[i].TargetDir, conf.TargetDir),
						UniqueNum: currentUniqueNum,
						Paused:    f.Paused,
					})

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
					fmt.Println("hit:", v.Title)

					hits = append(hits, Hit{
						Labels:    utils.FilterEmptyStrings([]string{"trfeed", currentFilter.Label}),
						Title:     v.Title,
						Resource:  v.Link,
						TargetDir: handleTargetDirs(currentFilter.TargetDir, conf.Feeds[i].TargetDir, conf.TargetDir),
						UniqueNum: currentUniqueNum,
						Paused:    currentFilter.Paused,
					})

					continue outer
				}
			}

			// (3) get all
			if conf.Feeds[i].GetAll {
				fmt.Println("hit:", v.Title)

				hits = append(hits, Hit{
					Labels:    utils.FilterEmptyStrings([]string{"trfeed", conf.Feeds[i].Label}),
					Title:     v.Title,
					Resource:  v.Link,
					TargetDir: handleTargetDirs(conf.Feeds[i].TargetDir, conf.TargetDir),
					UniqueNum: currentUniqueNum,
					Paused:    conf.Feeds[i].Paused,
				})
			}
		}
	}

	// back up:
	if err := utils.PutYAMLToFile(remnantsFile, &hits); err != nil {
		fmt.Println(err)
	}
}

func AddHits() {
	hitsLen := len(hits)
	for i := hitsLen - 1; i >= 0; i-- {
		err := addtorrent.AddTorrentWithRemote(
			conf.Host,
			conf.Auth,
			hits[i].Resource,
			hits[i].TargetDir,
			hits[i].Labels,
			hits[i].Paused,
		)

		if err == nil {
			fmt.Println("torrent added successfully:", hits[i].Title)
			hits.Remove(i)
		} else {
			fmt.Println("torrent could not be added:", hits[i].Title)
		}

		if hitsLen >= 10 {
			time.Sleep(time.Millisecond * 500)
		}
	}

	if err := utils.PutYAMLToFile(remnantsFile, &hits); err != nil {
		fmt.Println(err)
	}
}
