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
	Feeds     []Feed `yaml:"feeds"`
	Host      string `yaml:"host"`
	Auth      string `yaml:"auth"`
	TargetDir string `yaml:"target_dir"`
}
