package utils

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

func GetChildPids(pid int) int {
	scmd := fmt.Sprintf("pstree %d -p | awk -F\"[()]\" '{for(i=0;i<=NF;i++)if($i~/[0-9]+/)print $i}'", pid)
	cmd := exec.Command("bash", "-c", scmd)
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	//for i := 0; i < len(lines); i++ {
	//	fmt.Println(lines[i])
	//}
	atoi, _ := strconv.Atoi(lines[2])
	return atoi
}
