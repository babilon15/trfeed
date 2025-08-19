package scan

type Hit struct {
	Labels    []string `yaml:"labels"`
	Title     string   `yaml:"title"`
	Resource  string   `yaml:"resource"`
	TargetDir string   `yaml:"target_dir"`
	UniqueNum uint64   `yaml:"unique_num"`
	Paused    bool     `yaml:"paused"`
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
