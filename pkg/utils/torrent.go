package utils

import (
	"github.com/anacrolix/torrent/metainfo"
)

const (
	minTorrentSize int64 = 53687091200 // 50 gigabytes
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
		return minTorrentSize
	}
	return size
}
