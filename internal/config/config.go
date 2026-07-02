package config

import (
	"log"
	"strings"

	"github.com/babilon15/trfeed/pkg/utils"
)

type Filter struct {
	TargetDirs []string `yaml:"target_dirs"`
	RelPath    string   `yaml:"rel_path"`
	Include    string   `yaml:"include"`
	Exclude    string   `yaml:"exclude"`
	Label      string   `yaml:"label"`
	Literally  bool     `yaml:"literally"`
	Paused     bool     `yaml:"paused"`
	Disabled   bool     `yaml:"disabled"`
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
	TargetDirs       []string `yaml:"target_dirs"`
	RelPath          string   `yaml:"rel_path"`
	Url              string   `yaml:"url"`
	Label            string   `yaml:"label"`
	GetAll           bool     `yaml:"get_all"`
	Paused           bool     `yaml:"paused"`
	NoGlobalFilters  bool     `yaml:"no_global_filters"`
}

type Config struct {
	Feeds            []Feed   `yaml:"feeds"`
	Filters          []Filter `yaml:"filters"` // GLOBAL!
	TargetDirs       []string `yaml:"target_dirs"`
	RelPath          string   `yaml:"rel_path"`
	Host             string   `yaml:"host"`
	Auth             string   `yaml:"auth"`
	NoSpaceMarginGB  int64    `yaml:"no_space_margin_gb"`
	PausedIfNoSpace  bool     `yaml:"paused_if_no_space"`
	RandomTargetDirs bool     `yaml:"random_target_dirs"`
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

func IsFilterEmpty(f Filter) bool {
	return f.Include == ""
}
