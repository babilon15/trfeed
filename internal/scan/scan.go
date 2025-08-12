package scan

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

		err := os.MkdirAll(d, utils.DMode)
		if err != nil {
			log.Println(err)
		} else {
			return d
		}
	}
	return ""
}

type Scanner struct {
	Conf config.Config
	Hits Hits
}

func (s *Scanner) Init() {
	if err := utils.GetYAMLFromFile(configFile, &s.Conf); err != nil {
		log.Println(err)
	}

	if err := utils.GetYAMLFromFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}
}

func (s *Scanner) Save() {
	if s.Conf.ConfigOverwrite {
		if err := utils.PutYAMLToFile(configFile, &s.Conf); err != nil {
			log.Println(err)
		}
	}

	if err := utils.PutYAMLToFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}
}

func (s *Scanner) checkHit(item *feed.Item, feedIndex int) {
	if s.Hits.IndexByUniqueNum(item.GetUniqueNum()) != -1 {
		return
	}

	// (1) own filters
	for i := 0; i < len(s.Conf.Feeds[feedIndex].Filters); i++ {
		if s.Conf.Feeds[feedIndex].Filters[i].Check(item.Title) {
			log.Println("hit:", strconv.Quote(item.Title), "pub. date:", item.GetPubDate())

			s.Hits = append(s.Hits, Hit{
				Labels:    utils.FilterEmptyStrings([]string{s.Conf.Feeds[feedIndex].Filters[i].Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
				Title:     item.Title,
				Resource:  item.Link,
				TargetDir: handleTargetDirs(s.Conf.Feeds[feedIndex].Filters[i].TargetDir, s.Conf.Feeds[feedIndex].TargetDir, s.Conf.TargetDir),
				UniqueNum: item.GetUniqueNum(),
				Paused:    s.Conf.Feeds[feedIndex].Filters[i].Paused,
			})

			return
		}
	}

	// (2) filters via labels
	for i := 0; i < len(s.Conf.Feeds[feedIndex].FiltersViaLabels); i++ {
		filter := s.Conf.GetFilterByLabel(s.Conf.Feeds[feedIndex].FiltersViaLabels[i])

		if config.IsFilterEmpty(filter) {
			continue
		}

		if filter.Check(item.Title) {
			log.Println("hit:", strconv.Quote(item.Title), "pub. date:", item.GetPubDate())

			s.Hits = append(s.Hits, Hit{
				Labels:    utils.FilterEmptyStrings([]string{filter.Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
				Title:     item.Title,
				Resource:  item.Link,
				TargetDir: handleTargetDirs(filter.TargetDir, s.Conf.Feeds[feedIndex].TargetDir, s.Conf.TargetDir),
				UniqueNum: item.GetUniqueNum(),
				Paused:    filter.Paused,
			})

			return
		}
	}

	// (3) get all
	if s.Conf.Feeds[feedIndex].GetAll {
		log.Println("hit:", strconv.Quote(item.Title), "pub. date:", item.GetPubDate())

		s.Hits = append(s.Hits, Hit{
			Labels:    utils.FilterEmptyStrings([]string{s.Conf.Feeds[feedIndex].Label, "trfeed"}),
			Title:     item.Title,
			Resource:  item.Link,
			TargetDir: handleTargetDirs(s.Conf.Feeds[feedIndex].TargetDir, s.Conf.TargetDir),
			UniqueNum: item.GetUniqueNum(),
			Paused:    s.Conf.Feeds[feedIndex].Paused,
		})

		return
	}
}

func (s *Scanner) Run() {
	if len(s.Conf.Feeds) == 0 {
		log.Println("no feeds could be found in the config file:", configFile)
		return
	}

	for i := 0; i < len(s.Conf.Feeds); i++ {
		if !utils.IsValidURL(s.Conf.Feeds[i].Url) {
			log.Println("invalid URL:", strconv.Quote(s.Conf.Feeds[i].Url))
			continue
		}

		var currentFeed *feed.Feed

		if err := GetFeed(s.Conf.Feeds[i].Url, &currentFeed); err != nil {
			log.Println(err)
			continue
		}

		if len(currentFeed.Channel.Item) == 0 {
			log.Println("no items could be found in the feed", "url:", strconv.Quote(s.Conf.Feeds[i].Url))
			continue
		}

		currentFirstUniqueNum := currentFeed.Channel.Item[0].GetUniqueNum()

		lastItemFound := false

		for j := 0; j < len(currentFeed.Channel.Item); j++ {
			currentUniqueNum := currentFeed.Channel.Item[j].GetUniqueNum()
			if currentUniqueNum == s.Conf.Feeds[i].LastUniqueNum {
				lastItemFound = true
				break
			}

			s.checkHit(&currentFeed.Channel.Item[j], i)
		}

		if !lastItemFound && s.Conf.Feeds[i].LastUniqueNum != 0 {
			log.Println("the last checked item could not be found; it may have dropped from the feed;", "url:", strconv.Quote(s.Conf.Feeds[i].Url))
		}

		s.Conf.Feeds[i].LastUniqueNum = currentFirstUniqueNum

		time.Sleep(time.Millisecond * 250)
	}
}

func (s *Scanner) AddHits() {
	hitsLen := len(s.Hits)

	for i := hitsLen - 1; i >= 0; i-- {
		err := addtorrent.AddTorrentWithRemote(
			s.Conf.Host,
			s.Conf.Auth,
			s.Hits[i].Resource,
			s.Hits[i].TargetDir,
			s.Hits[i].Labels,
			s.Hits[i].Paused,
		)

		if err == nil {
			log.Println("torrent added successfully:", s.Hits[i].Title)
			s.Hits.Remove(i)
		} else {
			log.Println("torrent could not be added:", s.Hits[i].Title)
		}

		if hitsLen >= 5 {
			time.Sleep(time.Millisecond * 500)
		}
	}
}
