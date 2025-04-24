package prefixes

import (
	"fmt"
	"math"
)

const (
	precision    int     = 1
	prefixMaxNum int     = 9
	decBase      float64 = 1000.0
	binBase      float64 = 1024.0
)

func prefix(bytes float64, base float64) (float64, int) {
	var i int
	for bytes >= base && i < prefixMaxNum {
		bytes = bytes / base
		i++
	}
	return bytes, i
}

func decPrefix(bytes float64) (float64, string) {
	units := [9]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	size, prefixNum := prefix(bytes, decBase)
	return size, units[prefixNum]
}

func binPrefix(bytes float64) (float64, string) {
	units := [9]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	size, prefixNum := prefix(bytes, binBase)
	return size, units[prefixNum]
}

func roundFloat(num float64, p int) float64 {
	ratio := math.Pow(10, float64(p))
	return math.Round(num*ratio) / ratio
}

func GetPrefixSize(bytes int64) string {
	size, unit := binPrefix(float64(bytes))
	size = roundFloat(size, precision)
	return fmt.Sprintf("%g %s", size, unit)
}

func GetDecPrefixSize(bytes int64) string {
	size, unit := decPrefix(float64(bytes))
	size = roundFloat(size, precision)
	return fmt.Sprintf("%g %s", size, unit)
}
