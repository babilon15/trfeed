package diskusage

import (
	"syscall"
)

type DiskUsage struct {
	stat *syscall.Statfs_t
}

func (du *DiskUsage) Free() int64 {
	return int64(du.stat.Bfree) * du.stat.Bsize
}

func (du *DiskUsage) Available() int64 {
	return int64(du.stat.Bavail) * du.stat.Bsize
}

func (du *DiskUsage) Size() int64 {
	return int64(du.stat.Blocks) * du.stat.Bsize
}

func (du *DiskUsage) Used() int64 {
	return du.Size() - du.Free()
}

func GetDiskUsage(path string) *DiskUsage {
	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)
	return &DiskUsage{&stat}
}
