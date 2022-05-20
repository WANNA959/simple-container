package container

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"simple-container/pkg/cgroups"
	"simple-container/pkg/cgroups/subsystems"
	"simple-container/pkg/network"
	"simple-container/pkg/sqlite"
	"simple-container/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type ContainerInfo struct {
	Id          string `json:"id"`
	Pid         string `json:"pid"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
}

const (
	DefaultInfoLocation = "/var/run/simple-container"
	ConfigName          = "config.json"
	StateRunning        = "Running"
	StateStop           = "Stopped"
)

func RunWithCommand(tty bool, res *subsystems.ResourceConfig, containerName, net string) error {
	var netnsName string
	flag := false
	if net != "" {
		netnsName = net
		// pre-check netns, if not exist, create it
		if !network.NetnsExist(netnsName) {
			network.AddNetns(netnsName)
		}
		net = fmt.Sprintf("--net=/var/run/netns/%s", net)
	} else {
		flag = true
		netnsName = "netns" + strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
		net = fmt.Sprintf("--net")
		// via soft link
		//net = fmt.Sprintf("--net=/var/run/netns/%s", netnsName)
		//if !network.NetnsExist(netnsName) {
		//	network.AddNetns(netnsName)
		//}
	}
	scmd := fmt.Sprintf("unshare --ipc --user --uts %s --mount --root /root/cloud/centos/ --pid --mount-proc --fork bash", net)
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

	//time.Sleep(3 * time.Millisecond)
	// set uid & gid
	childPid := cmd.Process.Pid
	for {
		if err := utils.SetUserNsId(childPid); err != nil {
			log.Println(err)
			//return err
		} else {
			break
		}
	}

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
	cname, err := WriteContainerInfo(containerName, childPid)
	if err != nil {
		return err
	}

	// cgroups
	cgroupName := fmt.Sprintf("sc-cgroup-%d", childPid)
	cgroupManager := cgroups.NewCgroupManager(cgroupName)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(childPid)

	defer clean(netnsName, cname)
	cmd.Wait()
	return nil
}

func WriteContainerInfo(containerName string, pid int) (string, error) {
	id := strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
	createTime := time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02 15:04:05")

	if containerName == "" {
		containerName = id
	}
	containerInfo := &ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(pid),
		Name:        containerName,
		Command:     "unshare",
		CreatedTime: createTime,
		Status:      StateRunning,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Panicf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := filepath.Join(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Panicf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := filepath.Join(dirUrl, ConfigName)
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Fatalf("File write string error %v", err)
		return "", err
	}
	return containerName, nil
}

// use self metadata
func GetContainerInfo(file os.FileInfo) (*ContainerInfo, error) {
	containerName := file.Name()
	configFileDir := filepath.Join(DefaultInfoLocation, containerName, ConfigName)
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		log.Fatalf("Read file %s error %v", configFileDir, err)
		return nil, err
	}
	var containerInfo ContainerInfo

	// 将json文件信息反序列化成容器信息对象
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Fatalf("Json unmarshal error %v", err)
		return nil, err
	}

	return &containerInfo, nil
}

// no self metadata
func GetContainerInfoBack(pid string) *ContainerInfo {
	startTime, _ := utils.ProcessStartTime(pid)
	container := &ContainerInfo{
		Pid:         pid,
		Name:        "",
		Command:     utils.ProcessComm(pid),
		CreatedTime: startTime,
		Status:      utils.ProcessStat(pid),
	}
	return container
}

func RemoveContainerInfo(containerName string) {
	dirURL := filepath.Join(DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Fatalf("Remove file %s error %v", dirURL, err)
		return
	}
}

func clean(name, cname string) {
	network.DeleteDirNetns(name)
	network.DeleteNetns(name)
	RemoveContainerInfo(cname)
}

func cleanDb(ip string) {
	nm := sqlite.NetworkMgr{}
	nm.DeleteByBindIp(strings.SplitN(ip, "/", 2)[0])
}
