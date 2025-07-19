package scan

import (
	"fmt"

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
}

func Scan() {
	if err := utils.PutYAMLToFile(remnantsFile, &hits); err != nil {
		fmt.Println(err)
	}
}

func AddHits() {
	for i := 0; i < len(hits); i++ {
		fmt.Println(hits[i])
	}
}
