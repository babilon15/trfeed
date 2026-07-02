package torrent

import (
	"github.com/anacrolix/torrent/metainfo"
)

func GetInfo(path string) (metainfo.Info, error) {
	mi, miErr := metainfo.LoadFromFile(path)
	if miErr != nil {
		return metainfo.Info{}, miErr
	}

	return mi.UnmarshalInfo()
}

func GetTorrentSize(path string) int64 {
	info, _ := GetInfo(path)
	return info.TotalLength()
}

func IsMagnetLink(m string) bool {
	_, err := metainfo.ParseMagnetUri(m)
	return err == nil
}
