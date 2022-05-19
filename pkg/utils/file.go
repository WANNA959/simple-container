package utils

import "C"
import (
	"fmt"
	"io/ioutil"
	"os"
	"simple-container/pkg/container"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// only all files exist return true, other return false
func Exists(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			return false
		}
	}
	return true
}

func NotExists(files ...string) bool {
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			return false
		}
	}
	return true
}

func GetAllPid() []string {
	cpuPath := "/sys/fs/cgroup/cpu"
	memoryPath := "/sys/fs/cgroup/memory"
	files, _ := ioutil.ReadDir(cpuPath)
	pidsMap := make(map[string]bool)
	for _, file := range files {
		if strings.Contains(file.Name(), "sc-group") {
			pid := strings.ReplaceAll(file.Name(), "sc-group-", "")
			pidsMap[pid] = true
		}
	}
	files, _ = ioutil.ReadDir(memoryPath)
	for _, file := range files {
		if strings.Contains(file.Name(), "sc-group") {
			pid := strings.ReplaceAll(file.Name(), "sc-group-", "")
			pidsMap[pid] = true
		}
	}

	pids := make([]string, 0)
	for k, _ := range pidsMap {
		pids = append(pids, k)
	}
	return pids
}

func getUpdateTime() (int64, int64) {
	sys := syscall.Sysinfo_t{}
	syscall.Sysinfo(&sys)
	return time.Now().Unix() - sys.Uptime, int64(C.sysconf(C._SC_CLK_TCK))
}

func ProcessStartTime(pid string) (ts time.Time) {
	Uptime, scClkTck := getUpdateTime()
	buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%v/stat", pid))
	if err != nil {
		return time.Unix(0, 0)
	}
	if fields := strings.Fields(string(buf)); len(fields) > 22 {
		start, err := strconv.ParseInt(fields[21], 10, 0)
		if err == nil {
			if scClkTck > 0 {
				return time.Unix(Uptime+(start/scClkTck), 0)
			}
			return time.Unix(Uptime+(start/100), 0)
		}
	}
	return time.Unix(0, 0)
}

func GetContainerInfo(pid string) *container.ContainerInfo {
	container := &container.ContainerInfo{
		Pid:         pid,
		Name:        "",
		Command:     "sh",
		CreatedTime: ProcessStartTime(pid).String(),
		Status:      "Running",
	}
	return container
}
