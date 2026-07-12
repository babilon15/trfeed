package config

import (
	"log"
	"strings"

	"github.com/babilon15/trfeed/pkg/utils"
)

type Filter struct {
	TargetDirs []string `yaml:"target_dirs" json:"target_dirs"`
	RelPath    string   `yaml:"rel_path" json:"rel_path"`
	Include    string   `yaml:"include" json:"include"`
	Exclude    string   `yaml:"exclude" json:"exclude"`
	Label      string   `yaml:"label" json:"label"`
	Literally  bool     `yaml:"literally" json:"literally"`
	Pause      bool     `yaml:"pause" json:"pause"`
	Disabled   bool     `yaml:"disabled" json:"disabled"`
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
	Filters          []Filter `yaml:"filters" json:"filters"`
	FiltersViaLabels []string `yaml:"filters_via_labels" json:"filters"`
	TargetDirs       []string `yaml:"target_dirs" json:"filters"`
	RelPath          string   `yaml:"rel_path" json:"filters"`
	Url              string   `yaml:"url" json:"filters"`
	Label            string   `yaml:"label" json:"label"`
	GetAll           bool     `yaml:"get_all" json:"get_all"`
	Pause            bool     `yaml:"pause" json:"pause"`
	NoGlobalFilters  bool     `yaml:"no_global_filters" json:"no_global_filters"`
}

type Config struct {
	Feeds                []Feed   `yaml:"feeds" json:"feeds"`
	Filters              []Filter `yaml:"filters" json:"filters"` // GLOBAL!
	TargetDirs           []string `yaml:"target_dirs" json:"target_dirs"`
	RelPath              string   `yaml:"rel_path" json:"rel_path"`
	Host                 string   `yaml:"host" json:"host"`
	Auth                 string   `yaml:"auth" json:"auth"`
	TrRemotePath         string   `yaml:"tr_remote_path" json:"tr_remote_path"`
	NoSpaceMarginGB      int64    `yaml:"no_space_margin_gb" json:"no_space_margin_gb"`
	PauseIfNoSpace       bool     `yaml:"pause_if_no_space" json:"pause_if_no_space"`
	MaxFilesBeforePause  int      `yaml:"max_files_before_pause" json:"max_files_before_pause"`     // disable -> less than 1
	MaxSizeGBBeforePause int64    `yaml:"max_size_gb_before_pause" json:"max_size_gb_before_pause"` // disable -> less than 1
	RandomTargetDirs     bool     `yaml:"random_target_dirs" json:"random_target_dirs"`
	ReloadConfigFile     bool     `yaml:"reload_config_file" json:"reload_config_file"`
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
