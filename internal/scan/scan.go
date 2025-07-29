package scan

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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
	if err := utils.PutYAMLToFile(configFile, &s.Conf); err != nil {
		log.Println(err)
	}

	if err := utils.PutYAMLToFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}
}

func (s *Scanner) checkHit(item *feed.Item) {
	log.Println(item.Title)
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
			//***//
			s.checkHit(&currentFeed.Channel.Item[j])
			//***//
		}

		if !lastItemFound && s.Conf.Feeds[i].LastUniqueNum != 0 {
			log.Println("the last checked item could not be found; it may have dropped from the feed")
		}

		s.Conf.Feeds[i].LastUniqueNum = currentFirstUniqueNum
	}
}
