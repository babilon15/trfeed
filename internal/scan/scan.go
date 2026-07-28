package scan

import (
	"encoding/xml"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/babilon15/trfeed/internal/addtorrent"
	"github.com/babilon15/trfeed/internal/config"
	"github.com/babilon15/trfeed/pkg/diskusage"
	"github.com/babilon15/trfeed/pkg/feed"
	"github.com/babilon15/trfeed/pkg/torrent"
	"github.com/babilon15/trfeed/pkg/utils"
)

const (
	configFile       = "config.yaml"
	configFileJSON   = "config.json"
	remnantsFile     = "remnants.json"
	lastidsFile      = "lastids.json"
	torrentTargetDir = "trfeed"
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

func handleTargetDirs(dirs ...[]string) []string {
	for _, d := range dirs {
		if len(d) != 0 {
			c := utils.FilterEmptyStrings(d)
			if len(c) != 0 {
				return c
			}
		}
	}
	return []string{}
}

type Scanner struct {
	Conf          config.Config
	LastIDs       LastIDs
	Hits          Hits
	torrentTarget string
	useJSON       bool
}

func (s *Scanner) GetConfigFile() {
	if s.useJSON {
		if err := utils.GetJSONFromFile(configFileJSON, &s.Conf); err != nil {
			log.Println(err)
		}

		return
	}

	if err := utils.GetYAMLFromFile(configFile, &s.Conf); err != nil {
		log.Println(err)
	}
}

func (s *Scanner) Init(useJSON bool) {
	s.useJSON = useJSON

	s.GetConfigFile()

	if err := utils.GetJSONFromFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}

	if err := utils.GetJSONFromFile(lastidsFile, &s.LastIDs); err != nil {
		log.Println(err)
	}

	if s.Conf.NoSpaceMarginGB == 0 {
		s.Conf.NoSpaceMarginGB = 10
	}

	if s.torrentTarget == "" {
		s.torrentTarget = filepath.Join(os.TempDir(), torrentTargetDir)
		err := os.MkdirAll(s.torrentTarget, utils.DMode)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Scanner) Save() {
	if err := utils.PutJSONToFile(remnantsFile, &s.Hits); err != nil {
		log.Println(err)
	}

	if err := utils.PutJSONToFile(lastidsFile, &s.LastIDs); err != nil {
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
				Labels:     utils.FilterEmptyStrings([]string{s.Conf.Feeds[feedIndex].Filters[i].Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
				Title:      item.Title,
				Resource:   item.Link,
				TargetDirs: handleTargetDirs(s.Conf.Feeds[feedIndex].Filters[i].TargetDirs, s.Conf.Feeds[feedIndex].TargetDirs, s.Conf.TargetDirs),
				RelPath:    s.Conf.Feeds[feedIndex].Filters[i].RelPath,
				UniqueNum:  item.GetUniqueNum(),
				Pause:      s.Conf.Feeds[feedIndex].Filters[i].Pause,
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
				Labels:     utils.FilterEmptyStrings([]string{filter.Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
				Title:      item.Title,
				Resource:   item.Link,
				TargetDirs: handleTargetDirs(filter.TargetDirs, s.Conf.Feeds[feedIndex].TargetDirs, s.Conf.TargetDirs),
				RelPath:    filter.RelPath,
				UniqueNum:  item.GetUniqueNum(),
				Pause:      filter.Pause, // IMPORTANT!
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
					Labels:     utils.FilterEmptyStrings([]string{s.Conf.Filters[i].Label, s.Conf.Feeds[feedIndex].Label, "trfeed"}),
					Title:      item.Title,
					Resource:   item.Link,
					TargetDirs: handleTargetDirs(s.Conf.Filters[i].TargetDirs, s.Conf.Feeds[feedIndex].TargetDirs, s.Conf.TargetDirs),
					RelPath:    s.Conf.Filters[i].RelPath,
					UniqueNum:  item.GetUniqueNum(),
					Pause:      s.Conf.Filters[i].Pause,
				})

				return
			}
		}
	}

	// (4) GET ALL
	if s.Conf.Feeds[feedIndex].GetAll {
		log.Println("hit:", strconv.Quote(item.Title), "pub. date:", item.GetPubDate())

		s.Hits = append(s.Hits, Hit{
			Labels:     utils.FilterEmptyStrings([]string{s.Conf.Feeds[feedIndex].Label, "trfeed"}),
			Title:      item.Title,
			Resource:   item.Link,
			TargetDirs: handleTargetDirs(s.Conf.Feeds[feedIndex].TargetDirs, s.Conf.TargetDirs),
			RelPath:    s.Conf.Feeds[feedIndex].RelPath,
			UniqueNum:  item.GetUniqueNum(),
			Pause:      s.Conf.Feeds[feedIndex].Pause,
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
		if torrent.IsMagnetLink(s.Hits[i].Resource) {
			log.Println("magnet links are not currently supported")
			s.Hits.Remove(i)
			continue
		}

		dlPath, dlErr := DownloadTorrentFile(s.Hits[i].Resource, s.torrentTarget)
		if dlErr != nil {
			log.Println(dlErr)
		}

		torrentInfo, torrentInfoErr := torrent.GetTInfo(dlPath)
		if torrentInfoErr != nil {
			log.Println(torrentInfoErr)
		}

		if torrentInfo.Size == 0 {
			log.Println("torrent total size could not be determined:", strconv.Quote(s.Hits[i].Title))
			continue
		}

		targetDirs := slices.Clone(s.Hits[i].TargetDirs)

		if s.Conf.RandomTargetDirs {
			rand.Shuffle(len(targetDirs), func(i, j int) {
				targetDirs[i], targetDirs[j] = targetDirs[j], targetDirs[i]
			})
		}

		tdRes := ""
		for _, d := range targetDirs {
			if !filepath.IsAbs(d) {
				log.Println("this is not an absolute path:", strconv.Quote(d))
				continue
			}

			// usage := diskusage.GetDiskUsage(d)
			// -> It does not work properly if the path does not exist.

			usage := diskusage.GetDiskUsageIndefPath(d)

			if usage.Available() >= torrentInfo.Size+(s.Conf.NoSpaceMarginGB*1073741824) {
				tdRes = d
				break
			}
		}

		pause := s.Hits[i].Pause
		if tdRes == "" {
			if len(targetDirs) != 0 {
				tdRes = targetDirs[0]
			}

			if s.Conf.PauseIfNoSpace {
				pause = true
			}
		}

		if s.Conf.MaxFilesBeforePause > 0 {
			if torrentInfo.FilesNum > s.Conf.MaxFilesBeforePause {
				pause = true
			}
		}

		if s.Conf.MaxSizeGBBeforePause > 0 {
			if torrentInfo.Size > (s.Conf.MaxSizeGBBeforePause * 1073741824) {
				pause = true
			}
		}

		err := addtorrent.AddTorrentWithRemote(
			s.Conf.Host,
			s.Conf.Auth,
			dlPath,
			filepath.Join(tdRes, s.Hits[i].RelPath),
			s.Hits[i].Labels,
			pause,
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
