package config

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
	TargetMode    string   `yaml:"target_mode"` // [order], random, first
	TargetDirs    []string `yaml:"target_dirs"`
	Include       string   `yaml:"include"`
	Exclude       string   `yaml:"exclude"`
	Label         string   `yaml:"label"`
}

func (f *Filter) GetTargetMode() string {
	allowed := []string{"order", "random", "first"}
	if slices.Contains(f.TargetMode, allowed) {
		return f.TargetMode
	}
	return "order" // default
}

type Feed struct {
	GetAll     bool     `yaml:"get_all"`
	Paused     bool     `yaml:"paused"`
	LastNum    uint32   `yaml:"_last_num"` // FNV - 32-bit
	Url        string   `yaml:"url"`
	Label      string   `yaml:"label"`
	TargetMode string   `yaml:"target_mode"`
	TargetDirs []string `yaml:"target_dirs"`
	Filters    []Filter `yaml:"filters"`
}

type Config struct {
	NoLabels   bool     `yaml:"no_labels"`
	WaitSec    int      `yaml:"wait_sec"`
	TrEndpoint string   `yaml:"tr_endpoint"`
	TargetDirs []string `yaml:"target_dirs"`
	Feeds      []Feed   `yaml:"feeds"`
}
