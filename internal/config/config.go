package config

type Filter struct {
	Include   string `yaml:"include"`
	Exclude   string `yaml:"exclude"`
	Label     string `yaml:"label"`
	RelDir    string `yaml:"rel_dir"`
	TargetDir string `yaml:"target_dir"`
	Literally bool   `yaml:"literally"`
	Paused    bool   `yaml:"paused"`
}

type Feed struct {
	Filters          []Filter `yaml:"filters"`
	FiltersViaLabels []string `yaml:"filters_via_labels"`
	Url              string   `yaml:"url"`
	Label            string   `yaml:"label"`
	RelDir           string   `yaml:"rel_dir"`
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
	TargetDir string   `yaml:"target_dir"`
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
