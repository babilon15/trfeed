package scan

type Hit struct {
	Labels    []string `yaml:"labels"`
	Resource  string   `yaml:"resource"`
	TargetDir string   `yaml:"target_dir"`
	UniqueNum uint32   `yaml:"unique_num"`
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
