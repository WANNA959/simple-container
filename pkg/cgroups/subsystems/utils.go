package subsystems

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"simple-container/pkg/utils"
	"strings"
)

func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// 查找cgroup在文件系统中的绝对路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || autoCreate {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {

			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}

// cgroup related command
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

memory.limit_in_bytes=2097152 # 2MB
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
