package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"simple-container/pkg/cgroups"
	"simple-container/pkg/cgroups/subsystems"
	"simple-container/pkg/utils"
	"time"
)

type ContainerInfo struct {
	Pid         string `json:"pid"` //容器init进程在宿主机的PID
	Name        string `json:"name"`
	Command     string `json:"command"` //容器init的运行命令
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
}

func RunWithCommand(tty bool, res *subsystems.ResourceConfig, net string) error {
	scmd := fmt.Sprintf("unshare --ipc --user --uts --net=/var/run/netns/%s --mount --root /root/cloud/centos/ --pid --mount-proc --fork bash", net)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s\n", scmd)
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// no blocked(vs Run) return directly
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}

	time.Sleep(5 * time.Millisecond)
	// set uid & gid
	childPid := cmd.Process.Pid
	if err := utils.SetUserNsId(childPid); err != nil {
		log.Println(err)
		//return err
	}

	// cgroups
	cgroupName := fmt.Sprintf("sc-cgroup-%d", childPid)
	cgroupManager := cgroups.NewCgroupManager(cgroupName)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(childPid)

	cmd.Wait()
	return nil
}
