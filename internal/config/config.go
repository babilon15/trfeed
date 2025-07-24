package config

import (
	"fmt"
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
	include_words := strings.Fields(f.Include)
	exclude_words := strings.Fields(f.Exclude)

	if !f.Literally {
		for i := 0; i < len(include_words); i++ {
			include_words[i] = strings.ToLower(include_words[i])
			rd, err := utils.RemoveDiacritics(include_words[i])
			if err != nil {
				fmt.Println(err)
			}
			include_words[i] = rd
		}

		for i := 0; i < len(exclude_words); i++ {
			exclude_words[i] = strings.ToLower(exclude_words[i])
			rd, err := utils.RemoveDiacritics(exclude_words[i])
			if err != nil {
				fmt.Println(err)
			}
			exclude_words[i] = rd
		}

		title = strings.ToLower(title)
		title, _ = utils.RemoveDiacritics(title)
	}

	i_hit, e_hit := 0, 0

	for _, v := range include_words {
		if strings.Contains(title, v) {
			i_hit++
		}
	}

	for _, v := range exclude_words {
		if strings.Contains(title, v) {
			e_hit++
		}
	}

	return i_hit == len(include_words) && e_hit == 0
}

type Feed struct {
	Filters          []Filter `yaml:"filters"`
	FiltersViaLabels []string `yaml:"filters_via_labels"`
	Url              string   `yaml:"url"`
	Label            string   `yaml:"label"`
	TargetDir        string   `yaml:"target_dir"`
	LastUniqueNum    uint32   `yaml:"last_unique_num"`
	GetAll           bool     `yaml:"get_all"`
	Paused           bool     `yaml:"paused"`
}

type Config struct {
	Feeds     []Feed   `yaml:"feeds"`
	Filters   []Filter `yaml:"filters"`
	Host      string   `yaml:"host"`
	Auth      string   `yaml:"auth"`
	TargetDir string   `yaml:"target_dir"` // last case
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
