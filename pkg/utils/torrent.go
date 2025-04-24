package utils

import (
	"github.com/anacrolix/torrent/metainfo"
)

const (
	minTorrentSize int64 = 53687091200 // 50 gigabyte
)

// Visszaadja a megadott torrent fájl alapján (path) a teljes letöltési méretét a torrentnek:
func GetTorrentSize(path string) (int64, error) {
	meta, metaErr := metainfo.LoadFromFile(path)
	if metaErr != nil {
		return 0, metaErr
	}
	info, infoErr := meta.UnmarshalInfo()
	return info.Length, infoErr
}

// Burkoló függvény a fentihez.
// Ha nem állapítható meg a tényleges méret, akkor a fenti állandó értékét adja vissza.
func GetTorrentEstSize(path string) int64 {
	size, _ := GetTorrentSize(path)
	if size <= 0 {
		return minTorrentSize
	}
	return size
}
