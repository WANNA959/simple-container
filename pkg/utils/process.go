package utils

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
)

func CheckPidExist(pid int) bool {
	path := filepath.Join("/proc", strconv.FormatInt(int64(pid), 10))
	return Exists(path)
}

func SetUserNsId(pid int) error {
	pids := strconv.FormatInt(int64(pid), 10)
	scmd := fmt.Sprintf("echo '0 0 1' > /proc/%s/uid_map; echo '0 0 1' > /proc/%s/gid_map", pids, pids)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	return err
}
