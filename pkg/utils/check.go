package utils

import (
	"log"
	"strconv"

	"github.com/babilon15/trfeed/pkg/diskusage"
	"github.com/babilon15/trfeed/pkg/prefixes"
)

func CheckFreeSpace(path string, margins ...int64) bool {
	var margin int64
	for _, v := range margins {
		margin += v
	}

	usage := diskusage.GetDiskUsage(path)

	size := usage.Size()
	avai := usage.Available()

	if (avai - margin) <= 0 {
		log.Println("not enough space:", strconv.Quote(path),
			"size:", prefixes.GetPrefixSize(size),
			"available:", prefixes.GetPrefixSize(avai),
			"needed:", prefixes.GetPrefixSize(margin),
		)

		return false // nok
	}

	return true // ok
}
