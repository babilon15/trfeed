package scan

import (
	"fmt"

	"github.com/babilon15/trfeed/internal/addtorrent"
	"github.com/babilon15/trfeed/internal/config"
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
			fmt.Println("The client's default download directory is unknown. The host may be unavailable.")
		} else {
			conf.TargetDir = dd
			fmt.Println("New 'target_dir' directory from the client:", dd)
		}
	}
}

func Scan() {
	feedsLen := len(conf.Feeds)

	if feedsLen == 0 {
		fmt.Println("No feeds could be found in the config file:", configFile)
	}

	for i := 0; i < feedsLen; i++ {
		// x//
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
