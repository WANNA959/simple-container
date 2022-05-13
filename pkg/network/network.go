package network

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"os/exec"
	"strings"
)

const (
	Veth = "veth"
)

func AddNetns(name string) error {
	scmd := fmt.Sprintf("ip netns add net %s", name)
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
	// assign ip subnet
	scmd := fmt.Sprintf("ip netns exec %s ip addr add %s dev %s", name, subnet, iface)
	cmd := exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// interface up
	scmd = fmt.Sprintf("ip netns exec %s ip link set %s up", name, iface)
	cmd = exec.Command("bash", "-c", scmd)
	log.Printf("exec command: %s", scmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
