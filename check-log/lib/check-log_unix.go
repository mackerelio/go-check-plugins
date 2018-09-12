// +build !windows

package checklog

import (
	"os"
	"syscall"
)

func detectInode(fi os.FileInfo) uint {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return uint(stat.Ino)
	}
	return 0
}
