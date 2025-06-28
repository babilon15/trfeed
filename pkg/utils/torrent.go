package utils

import (
	"net/url"

	"github.com/anacrolix/torrent/metainfo"
)

const (
	MinTorrentSize int64 = 53687091200 // 50 gigabytes
)

func GetTorrentSize(path string) (int64, error) {
	meta, metaErr := metainfo.LoadFromFile(path)
	if metaErr != nil {
		return 0, metaErr
	}
	info, infoErr := meta.UnmarshalInfo()
	return info.Length, infoErr
}

func GetTorrentEstSize(path string) int64 {
	size, _ := GetTorrentSize(path)
	if size <= 0 {
		return MinTorrentSize
	}
	return size
}

// If the link is a magnet link, we do not check the size of the torrent. (MinTorrentSize)
func IsMagnetLink(link string) bool {
	u, err := url.ParseRequestURI(link)
	return err == nil && u.Scheme == "magnet"
}
