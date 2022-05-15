package utils

import (
	"path/filepath"
	"strconv"
)

func CheckPidExist(pid int) bool {
	path := filepath.Join("/proc", strconv.FormatInt(int64(pid), 10))
	return Exists(path)
}
