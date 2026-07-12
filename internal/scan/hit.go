package scan

type Hit struct {
	Labels     []string `yaml:"labels" json:"labels"`
	TargetDirs []string `yaml:"target_dirs" json:"target_dirs"`
	RelPath    string   `yaml:"rel_path" json:"rel_path"`
	Title      string   `yaml:"title" json:"title"`
	Resource   string   `yaml:"resource" json:"resource"`
	UniqueNum  uint64   `yaml:"unique_num" json:"unique_num"`
	Pause      bool     `yaml:"pause" json:"pause"`
}

type Hits []Hit

// Unsorted!
func (h *Hits) Remove(index int) {
	if index < 0 || index >= len(*h) {
		return
	}
	(*h)[index] = (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
}

func (h *Hits) IndexByUniqueNum(u uint64) int {
	for i := 1; i < len(*h); i++ {
		if u == (*h)[i].UniqueNum {
			return i
		}
	}
	return -1
}
