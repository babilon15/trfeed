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
	lastidsFile  = "lastids.yaml"
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
	Conf    config.Config
	LastIDs LastIDs
	Hits    Hits
}

func (s *Scanner) Init() {
	if err := utils.GetYAMLFromFile(configFile, &s.Conf); err != nil {
		log.Println(err)
	}

	if err := utils.GetYAMLFromFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}

	if err := utils.GetYAMLFromFile(lastidsFile, &s.LastIDs); err != nil {
		log.Println(err)
	}

	if s.Conf.NoSpaceMarginGB == 0 {
		s.Conf.NoSpaceMarginGB = 50
	}
}

func (s *Scanner) Save() {
	if err := utils.PutYAMLToFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}

	if err := utils.PutYAMLToFile(lastidsFile, &s.LastIDs); err != nil {
		log.Println(err)
	}
}

func (s *Scanner) checkHit(item *feed.Item, feedIndex int, noGlobalFilters bool) {
	if s.Hits.IndexByUniqueNum(item.GetUniqueNum()) != -1 {
		return
	}

	// (1) OWN FILTERS
	for i := 0; i < len(s.Conf.Feeds[feedIndex].Filters); i++ {
		if s.Conf.Feeds[feedIndex].Filters[i].Disabled {
			continue
		}

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

	// (2) FILTERS VIA LABELS
	for i := 0; i < len(s.Conf.Feeds[feedIndex].FiltersViaLabels); i++ {
		filter := s.Conf.GetFilterByLabel(s.Conf.Feeds[feedIndex].FiltersViaLabels[i])

		if config.IsFilterEmpty(filter) || filter.Disabled {
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
				Paused:    filter.Paused, // IMPORTANT!
			})

			return
		}
	}

	// (3) GLOBAL
	if !noGlobalFilters {
		for i := 0; i < len(s.Conf.Filters); i++ {
			if s.Conf.Filters[i].Disabled {
				continue
			}

			if s.Conf.Filters[i].Check(item.Title) {
				log.Println("hit:", strconv.Quote(item.Title), "pub. date:", item.GetPubDate())

				s.Hits = append(s.Hits, Hit{
					Labels:    utils.FilterEmptyStrings([]string{s.Conf.Filters[i].Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
					Title:     item.Title,
					Resource:  item.Link,
					TargetDir: handleTargetDirs(s.Conf.Filters[i].TargetDir, s.Conf.Feeds[feedIndex].TargetDir, s.Conf.TargetDir),
					UniqueNum: item.GetUniqueNum(),
					Paused:    s.Conf.Filters[i].Paused,
				})

				return
			}
		}
	}

	// (4) GET ALL
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

		lastID := s.LastIDs.GetLastIDByUrl(s.Conf.Feeds[i].Url)

		for j := 0; j < len(currentFeed.Channel.Item); j++ {
			currentUniqueNum := currentFeed.Channel.Item[j].GetUniqueNum()
			if currentUniqueNum == lastID {
				lastItemFound = true
				break
			}

			s.checkHit(&currentFeed.Channel.Item[j], i, s.Conf.Feeds[i].NoGlobalFilters)
		}

		if !lastItemFound && lastID != 0 {
			log.Println("the last checked item could not be found; it may have dropped from the feed;", "url:", strconv.Quote(s.Conf.Feeds[i].Url))
		}

		s.LastIDs.SetLastIDByUrl(s.Conf.Feeds[i].Url, currentFirstUniqueNum)

		time.Sleep(time.Millisecond * 250)
	}
}

func (s *Scanner) AddHits() {
	hitsLen := len(s.Hits)

	for i := hitsLen - 1; i >= 0; i-- {
		if s.Hits[i].TargetDir != "" {
			if s.Conf.PausedIfNoSpace && !utils.CheckFreeSpace(s.Hits[i].TargetDir, s.Conf.NoSpaceMarginGB*1073741824) {
				s.Hits[i].Paused = true
			}
		}

		err := addtorrent.AddTorrentWithRemote(
			s.Conf.Host,
			s.Conf.Auth,
			s.Hits[i].Resource,
			s.Hits[i].TargetDir,
			s.Hits[i].Labels,
			s.Hits[i].Paused,
		)

		if err == nil {
			log.Println("torrent added successfully:", strconv.Quote(s.Hits[i].Title))
			s.Hits.Remove(i)
		} else {
			log.Println("torrent could not be added:", strconv.Quote(s.Hits[i].Title))
		}

		if hitsLen >= 5 {
			time.Sleep(time.Millisecond * 500)
		}
	}
}
