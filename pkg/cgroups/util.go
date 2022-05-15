package cgroups

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"simple-container/pkg/utils"
	"strings"
)

func DeleteCgroup(name string) error {
	scmd := fmt.Sprintf("cgdelete %s", name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func CreateCgroup(name string) error {
	scmd := fmt.Sprintf("cgcreate -g %s", name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func FindCgroupPath(name string) (string, error) {
	splitN := strings.SplitN(name, ":", 2)
	cType := splitN[0]
	cName := splitN[1]
	path := filepath.Join("/sys/fs/cgroup/", cType, cName)
	if utils.Exists(path) {
		return path, nil
	}
	return "", os.ErrNotExist
}

/*
cpu.shares=512
cpu.cfs_quota_us=10000
cpu.cfs_period_us=100000

memory.limit_in_bytes=2097152
memory.swappiness=0
memory.oom_control=1 # oom 1=wait 0=kill
*/
func SetCgroupLimit(limitName string, share string, name string) error {
	scmd := fmt.Sprintf("cgset -r %s=%s %s", limitName, share, name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
