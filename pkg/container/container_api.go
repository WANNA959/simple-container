package container

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"os/exec"
	"simple-container/pkg/cgroups"
	"simple-container/pkg/cgroups/subsystems"
	"simple-container/pkg/network"
	"simple-container/pkg/sqlite"
	"strings"
	"syscall"
)

func Run(tty bool, cmdArrays []string, res *subsystems.ResourceConfig, volume, containerName, imageName, net, comm string) error {
	cid := strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
	if containerName == "" {
		containerName = cid
	}

	var netnsName string
	flag := false
	if net != "" {
		netnsName = net
		// pre-check netns, if not exist, create it
		if !network.NetnsExist(netnsName) {
			network.AddNetns(netnsName)
		}
	} else {
		flag = true
		netnsName = "netns" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
	}

	parent, writePipe := NewParentProcess(tty, "/root/cloud/centos/", containerName)
	if parent == nil {
		log.Fatalln("New parent process error")
	}
	if err := parent.Start(); err != nil {
		log.Fatalln(err)
	}
	childPid := parent.Process.Pid

	//for {
	//	if err := utils.SetUserNsId(childPid); err != nil {
	//		log.Println(err)
	//		//return err
	//	} else {
	//		break
	//	}
	//}

	// connect to default docker0 bridge
	if flag {
		oldns := fmt.Sprintf("/proc/%d/ns/net", childPid)
		newns := fmt.Sprintf("/var/run/netns/%s", netnsName)
		if err := os.Symlink(oldns, newns); err != nil {
			return err
		}
		// create veth pair
		vethPairs, err := network.CraeteVethPair("", "")
		if err != nil {
			return err
		}
		cm := sqlite.ContainerMgr{}
		cm.Insert(sqlite.ContainerMgr{
			Pid:  netnsName,
			Veth: vethPairs[0],
		})

		// add veth to netns
		next, err := network.GetNextIp(network.DefaultBridgeSubnet)
		defer cleanDb(next)
		log.Printf("next ip:%s", next)
		if err != nil {
			return err
		}

		if err := network.AssignIpAndUp(netnsName, next, vethPairs[0]); err != nil {
			return err
		}

		// add veth to master bridge
		if err := network.AddVeth2BridgeNic(vethPairs[1], network.DefaultDocker0Bridge); err != nil {
			return err
		}
	}

	err := WriteContainerInfo(cid, containerName, imageName, volume, comm, childPid)
	if err != nil {
		return err
	}

	cgroupManager := cgroups.NewCgroupManager("simple-container")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArrays, writePipe)
	defer clean(netnsName, containerName)
	parent.Wait()
	os.Exit(0)
	return nil
}

// call init command here
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Printf("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func NewParentProcess(tty bool, rootfs string, containerName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Fatalf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Fatalf("NewParentProcess mkdir %s error %v", dirURL, err)
		return nil, nil
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = rootfs
	return cmd, writePipe
}

// 生成匿名管道
func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
