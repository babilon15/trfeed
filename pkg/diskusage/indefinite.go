package diskusage

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	mountsFile = "/proc/mounts"
)

func GetAbsPathDepth(path string) int {
	cPath, err := filepath.Abs(path)
	if err != nil {
		return 0 // I'm lazy...
	}

	if cPath == "/" {
		return 1
	}

	p := strings.FieldsFunc(cPath, func(r rune) bool {
		return r == '/'
	})

	return len(p) + 1
}

func GetMountpoint(path string) (string, error) {
	read, readErr := os.ReadFile(mountsFile)
	if readErr != nil {
		return path, readErr
	}

	lines := strings.Split(string(read), "\n")

	for _, l := range lines {
		fields := strings.Fields(l)

		if len(fields) != 6 || fields[1] == "/" {
			continue
		}

		if strings.HasPrefix(path, fields[1]) {
			return fields[1], nil
		}
	}

	return "/", nil
}

func GetDiskUsageIndefPath(path string) *DiskUsage {
	m, err := GetMountpoint(path)
	if err != nil {
		log.Println(err)
		return &DiskUsage{}
	}
	return GetDiskUsage(m)
}
