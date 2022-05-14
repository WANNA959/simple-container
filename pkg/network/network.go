package network

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"os/exec"
	"strings"
)

const (
	Veth                = "veth"
	DefaultMasterBridge = "master-br0"
)

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

func AddVeth2MasterNic(iface string) error {
	// set veth to netns
	scmd := fmt.Sprintf("ip link set dev %s master %s; ip link set dev %s up", iface, DefaultMasterBridge, iface)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func judgeNicExsit(masterBridge string) (bool, error) {
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

func GetBridgeSubnet(subnet string) (string, error) {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", err
	}
	last := int(ipNet.IP[3]) + 1
	ipNet.IP[3] = byte(last)
	splits := strings.SplitN(subnet, "/", 2)
	if splits[0] == ipNet.IP.String() {
		return "", errors.New(fmt.Sprintf("%s is reserved for bridge", ipNet.IP.String()))
	}
	return fmt.Sprintf("%s/%s", ipNet.IP.String(), splits[1]), nil
}

func GenerateBridgeOrSkip(subnet string) error {
	bridgeSubnet, err := GetBridgeSubnet(subnet)
	if err != nil {
		return err
	}
	exsit, err := judgeNicExsit(DefaultMasterBridge)
	if err != nil {
		return err
	}
	if !exsit {
		scmd := fmt.Sprintf("ip link add %s type bridge; ip addr add %s dev %s;  ip link set dev %s up", DefaultMasterBridge, bridgeSubnet, DefaultMasterBridge, DefaultMasterBridge)
		cmd := exec.Command("bash", "-c", scmd)
		log.Printf("exec command: %s", scmd)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
