package network

import (
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/pkg/network"
	"simple-container/pkg/sqlite"
	"text/template"
)

var connectBridgeTemplate = template.Must(template.New("simple-container controller connectBridge").Parse(`
------------------------------------------------
simple-container controller:
    master bridge name: {{.BridgeName}}
    master bridge subnet: {{.BridgeSubnet}}
    netns name: {{.Name}}
    netns subnet: {{.Subnet}}
------------------------------------------------
`))

func NewConnectBridgeCommand() *cli.Command {
	return &cli.Command{
		Name:      "connect-bridge",
		Usage:     "connect to host bridge",
		UsageText: "scadm [global options] connect-bridge [options]",
		Action:    connectBridge,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of netns",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "subnet",
				Usage:       "subnet of netns",
				Destination: &subnet,
			},
		},
	}
}

func connectBridge(ctx *cli.Context) error {

	// pre-check master bridge
	if err := network.GenerateBridgeOrSkip(subnet, network.DefaultMasterBridge); err != nil {
		return err
	}

	// pre-check netns, if not exist, create it
	if !network.NetnsExist(name) {
		network.AddNetns(name)
	}

	// create veth pair
	vethPairs, err := network.CraeteVethPair("", "")
	if err != nil {
		return err
	}

	// add veth to netns
	if err := network.AssignIpAndUp(name, subnet, vethPairs[0]); err != nil {
		return err
	}
	cm := sqlite.ContainerMgr{}
	cm.Insert(sqlite.ContainerMgr{
		Pid:  name,
		Veth: vethPairs[0],
	})

	// add veth to master bridge
	if err := network.AddVeth2BridgeNic(vethPairs[1], network.DefaultMasterBridge); err != nil {
		return err
	}

	bridgeSubnet, err := network.GetBridgeSubnet(subnet, false)
	if err != nil {
		return err
	}
	data := struct {
		BridgeName   string
		BridgeSubnet string
		Name         string
		Subnet       string
	}{
		BridgeName:   network.DefaultMasterBridge,
		BridgeSubnet: bridgeSubnet,
		Name:         name,
		Subnet:       subnet,
	}

	connectBridgeTemplate.Execute(os.Stdout, &data)

	return nil
}
