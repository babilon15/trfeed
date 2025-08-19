package config

import (
	"log"
	"strings"

	"github.com/babilon15/trfeed/pkg/utils"
)

type Filter struct {
	Include   string `yaml:"include"`
	Exclude   string `yaml:"exclude"`
	Label     string `yaml:"label"`
	TargetDir string `yaml:"target_dir"`
	Literally bool   `yaml:"literally"`
	Paused    bool   `yaml:"paused"`
}

func (f *Filter) Check(title string) bool {
	includeWords := strings.Fields(f.Include)
	excludeWords := strings.Fields(f.Exclude)

	if !f.Literally {
		for i := 0; i < len(includeWords); i++ {
			includeWords[i] = strings.ToLower(includeWords[i])
			rd, err := utils.RemoveDiacritics(includeWords[i])
			if err != nil {
				log.Println(err)
			}
			includeWords[i] = rd
		}

		for i := 0; i < len(excludeWords); i++ {
			excludeWords[i] = strings.ToLower(excludeWords[i])
			rd, err := utils.RemoveDiacritics(excludeWords[i])
			if err != nil {
				log.Println(err)
			}
			excludeWords[i] = rd
		}

		title = strings.ToLower(title)
		title, _ = utils.RemoveDiacritics(title)
	}

	iHit, eHit := 0, 0

	for _, v := range includeWords {
		if strings.Contains(title, v) {
			iHit++
		}
	}

	for _, v := range excludeWords {
		if strings.Contains(title, v) {
			eHit++
		}
	}

	return iHit == len(includeWords) && eHit == 0
}

type Feed struct {
	Filters          []Filter `yaml:"filters"`
	FiltersViaLabels []string `yaml:"filters_via_labels"`
	Url              string   `yaml:"url"`
	Label            string   `yaml:"label"`
	TargetDir        string   `yaml:"target_dir"`
	LastUniqueNum    uint64   `yaml:"last_unique_num"`
	GetAll           bool     `yaml:"get_all"`
	Paused           bool     `yaml:"paused"`
}

type Config struct {
	Feeds           []Feed   `yaml:"feeds"`
	Filters         []Filter `yaml:"filters"`
	Host            string   `yaml:"host"`
	Auth            string   `yaml:"auth"`
	TargetDir       string   `yaml:"target_dir"`
	ConfigOverwrite bool     `yaml:"config_overwrite"`
	PausedIfNoSpace bool     `yaml:"paused_if_no_space"`
}

func (c *Config) GetFilterByLabel(label string) Filter {
	for _, f := range c.Filters {
		if f.Label == label {
			return f
		}
	}

	for _, v := range c.Feeds {
		for _, f := range v.Filters {
			if f.Label == label {
				return f
			}
		}
	}

	return Filter{}
}

func IsFilterEmpty(f Filter) bool { return f == Filter{} }
