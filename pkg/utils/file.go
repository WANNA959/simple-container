package utils

import (
	"C"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
		if strings.Contains(file.Name(), "sc-cgroup") {
			pid := strings.ReplaceAll(file.Name(), "sc-cgroup-", "")
			pidsMap[pid] = true
		}
	}
	files, _ = ioutil.ReadDir(memoryPath)
	for _, file := range files {
		if strings.Contains(file.Name(), "sc-cgroup") {
			pid := strings.ReplaceAll(file.Name(), "sc-cgroup-", "")
			pidsMap[pid] = true
		}
	}

	pids := make([]string, 0)
	for k, _ := range pidsMap {
		pids = append(pids, k)
	}
	log.Printf("pids:%+v", pids)
	return pids
}

func ProcessStartTime(pid string) (string, error) {
	stat, err := os.Lstat(fmt.Sprintf("/proc/%v", pid))
	if err != nil {
		return "-1", err
	}

	// ?
	unix := time.Unix(stat.ModTime().Unix(), 0).Add(8 * time.Hour)
	return unix.Format("2006-01-02 15:04:05"), nil
}

func ProcessStat(pid string) string {
	buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%v/stat", pid))
	if err != nil {
		return "unknown"
	}
	stateMap := map[string]string{
		"R": "Running",
		"S": "Running",
		"D": "Stopped",
		"T": "Stopped",
		"Z": "Stopped",
		"X": "Stopped",
	}
	if fields := strings.Fields(string(buf)); len(fields) > 22 {
		return stateMap[fields[2]]
	}
	return "unknown"
}

func ProcessComm(pid string) string {
	buf, err := ioutil.ReadFile(fmt.Sprintf("/proc/%v/comm", pid))
	if err != nil {
		return "unknown"
	}
	return strings.ReplaceAll(string(buf), "\n", "")
}
