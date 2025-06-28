package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/babilon15/trfeed/internal/config"
	"github.com/babilon15/trfeed/internal/places"
	"github.com/babilon15/trfeed/internal/scan"
	"github.com/babilon15/trfeed/pkg/utils"
)

func main() {
	log.SetFlags(0)
	places.Checks()

	var conf config.Config

	if _, err := os.Stat(places.ConfigFile); errors.Is(err, os.ErrNotExist) {
		// It probably does not exist.
		utils.PutYAMLToFile(places.ConfigFile, &config.Config{})
	} else {
		// It exists.
		utils.GetYAMLFromFile(places.ConfigFile, &conf)
	}

	conf.Checks()

	for {
		scan.Scan()
		time.Sleep(time.Duration(conf.WaitSec))
	}
}
