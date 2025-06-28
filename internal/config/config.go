package config

import (
	"github.com/babilon15/trfeed/pkg/utils"
)

const (
	WAITSEC_MIN = 10   // Because we don't want to overload the servers.
	WAITSEC_MAX = 3600 // To prevent missing feed items.
)

// 'target_dirs' order:
// --------------------
// 1. Filter
// 2. Feed
// 3. Program default
// 4. Client default

type Filter struct {
	CaseSensitive bool     `yaml:"case_sensitive"`
	Diacritics    bool     `yaml:"diacritics"`
	Paused        bool     `yaml:"paused"`
	TargetDirs    []string `yaml:"target_dirs"`
	Include       string   `yaml:"include"`
	Exclude       string   `yaml:"exclude"`
	Label         string   `yaml:"label"`
}

type Feed struct {
	GetAll     bool     `yaml:"get_all"`
	Paused     bool     `yaml:"paused"`
	LastNum    uint32   `yaml:"_last_num"` // FNV - 32-bit
	Url        string   `yaml:"url"`
	Label      string   `yaml:"label"`
	TargetDirs []string `yaml:"target_dirs"`
	Filters    []Filter `yaml:"filters"`
}

type Config struct {
	NoLabels         bool     `yaml:"no_labels"`
	NoFreeSpaceCheck bool     `yaml:"no_free_space_check"`
	WaitSec          int      `yaml:"wait_sec"`
	TrEndpoint       string   `yaml:"tr_endpoint"`
	TargetDirs       []string `yaml:"target_dirs"`
	Feeds            []Feed   `yaml:"feeds"`
}

func (c *Config) Checks() {
	c.WaitSec = utils.Clamp(c.WaitSec, WAITSEC_MIN, WAITSEC_MAX)
}
