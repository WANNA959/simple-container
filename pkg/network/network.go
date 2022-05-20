package network

import (
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"os/exec"
	"path/filepath"
	"simple-container/pkg/sqlite"
	"simple-container/pkg/utils"
	"strconv"
	"strings"
)

const (
	Veth                 = "veth"
	DefaultMasterBridge  = "master-br0"
	DefaultDocker0Bridge = "sc-br0"
	DefaultBridgeSubnet  = "10.99.0.1/24"
)

var PoolFullErr = errors.New("IP Pool Full")

func NetnsExist(name string) bool {
	netnsPath := filepath.Join("/var/run/netns/", name)
	return utils.Exists(netnsPath)
}

func AddNetns(name string) error {
	scmd := fmt.Sprintf("ip netns add %s", name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func DeleteDirNetns(name string) error {
	scmd := fmt.Sprintf("rm -rf /var/run/netns/%s", name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func DeleteNetns(name string) error {
	scmd := fmt.Sprintf("ip netns delete %s", name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func randomVethPeer() []string {
	veth1 := strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
	veth2 := strings.ReplaceAll(uuid.NewV4().String(), "-", "")[:8]
	return []string{Veth + veth1, Veth + veth2}
}

func CraeteVethPair(veth1, veth2 string) ([]string, error) {
	// not assign, self generate
	if veth1 == "" && veth2 == "" {
		peers := randomVethPeer()
		veth1 = peers[0]
		veth2 = peers[1]
	}

	scmd := fmt.Sprintf("ip link add %s type veth peer name %s", veth1, veth2)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return []string{veth1, veth2}, nil
}

func AssignIpAndUp(name, subnet, iface string) error {
	// set veth to netns
	scmd := fmt.Sprintf("ip link set %s netns %s", iface, name)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// assign ip subnet
	scmd = fmt.Sprintf("ip netns exec %s ip addr add %s dev %s", name, subnet, iface)
	cmd = exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// interface up
	scmd = fmt.Sprintf("ip netns exec %s ip link set lo up; ip netns exec %s ip link set %s up", name, name, iface)
	cmd = exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func AddVeth2BridgeNic(iface, bridge string) error {
	// set veth to netns
	scmd := fmt.Sprintf("ip link set dev %s master %s; ip link set dev %s up", iface, bridge, iface)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func JudgeNicExsit(masterBridge string) (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for i := 0; i < len(interfaces); i++ {
		if interfaces[i].Name == masterBridge {
			return true, nil
		}
	}
	return false, nil
}

func GetBridgeSubnet(subnet string, isBridge bool) (string, error) {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", err
	}
	last := int(ipNet.IP[3]) + 1
	ipNet.IP[3] = byte(last)
	splits := strings.SplitN(subnet, "/", 2)
	if !isBridge && splits[0] == ipNet.IP.String() {
		return "", errors.New(fmt.Sprintf("%s is reserved for bridge", ipNet.IP.String()))
	}
	return fmt.Sprintf("%s/%s", ipNet.IP.String(), splits[1]), nil
}

func GenerateBridgeOrSkip(subnet, nicName string) error {
	bridgeSubnet, err := GetBridgeSubnet(subnet, true)
	if err != nil {
		return err
	}
	exsit, err := JudgeNicExsit(nicName)
	if err != nil {
		return err
	}
	if !exsit {
		scmd := fmt.Sprintf("ip link add %s type bridge; ip addr add %s dev %s;  ip link set dev %s up", nicName, bridgeSubnet, nicName, nicName)
		cmd := exec.Command("bash", "-c", scmd)
		log.Printf("exec command: %s", scmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNextIp(subnet string) (string, error) {
	info := strings.SplitN(subnet, "/", 2)
	ip := info[0]
	mask := info[1]
	nm := sqlite.NetworkMgr{}
	subnets, err := nm.QueryBySubnet(subnet)
	if err != nil && err != sqlite3.ErrNotFound {
		log.Fatalln(err)
		return "", err
	}
	ipMap := make(map[string]bool)
	for _, item := range subnets {
		splits := strings.SplitN(item.BindIp, ".", 4)
		ipMap[splits[3]] = true
	}
	for i := 2; i < 254; i++ {
		stri := strconv.FormatInt(int64(i), 10)
		if !ipMap[stri] {
			splits := strings.SplitN(ip, ".", 4)
			splits[3] = stri
			newip := strings.Join(splits, ".")
			nm := sqlite.NetworkMgr{}

			nm.Insert(sqlite.NetworkMgr{
				Subnet: subnet,
				BindIp: newip,
			})
			return newip + "/" + mask, nil
		}
	}
	return "", PoolFullErr
}
