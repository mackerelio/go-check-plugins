package checklog

import (
	"os"
)

func detectInode(_ os.FileInfo) uint {
	return 0
}
