package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"simple-container/pkg/cgroups"
	"simple-container/pkg/cgroups/subsystems"
)

func RunWithCommand(tty bool, res *subsystems.ResourceConfig) error {
	scmd := fmt.Sprintf("unshare --ipc --user --uts --net=/var/run/netns/netns3 --mount --root /root/cloud/centos/ --pid --mount-proc --fork bash")
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s\n", scmd)
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// no blocked(vs Run) return directly
	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	// cgroups
	cgroupManager := cgroups.NewCgroupManager("sc-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(cmd.Process.Pid)

	cmd.Wait()
	return nil
}
