package scan

type Pair struct {
	Url           string `yaml:"url"`
	LastUniqueNum uint64 `yaml:"last_unique_num"`
}

type LastIDs struct {
	Pairs []Pair `yaml:"last_ids"`
}

func (l *LastIDs) GetLastIDByUrl(url string) uint64 {
	for i := 0; i < len(l.Pairs); i++ {
		if l.Pairs[i].Url == url {
			return l.Pairs[i].LastUniqueNum
		}
	}
	return 0
}

func (l *LastIDs) SetLastIDByUrl(url string, uniqueNum uint64) {
	for i := 0; i < len(l.Pairs); i++ {
		if l.Pairs[i].Url == url {
			l.Pairs[i].LastUniqueNum = uniqueNum
			return
		}
	}

	l.Pairs = append(l.Pairs, Pair{Url: url, LastUniqueNum: uniqueNum})
}
