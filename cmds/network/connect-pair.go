package network

import (
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/pkg/network"
	"strings"
	"text/template"
)

// veth part mode

var connectPairTemplate = template.Must(template.New("simple-container controller connectPair").Parse(`
------------------------------------------------
simple-container controller:
    netns1
        name: {{.Name1}}
        subnet: {{.Subnet1}}
        veth name: {{.Veth1}}
    netns2
        name: {{.Name2}}
        subnet: {{.Subnet2}}
        veth name: {{.Veth2}}
------------------------------------------------
`))

var netns string
var subnets string

func NewConnectPairCommand() *cli.Command {
	return &cli.Command{
		Name:      "connect-pair",
		Usage:     "connect two netns with veth pair",
		UsageText: "scadm [global options] connect-pair [options]",
		Action:    connectPair,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "netns",
				Usage:       "names of netns",
				Destination: &netns,
			},
			&cli.StringFlag{
				Name:        "subnets",
				Usage:       "subnets of netns",
				Destination: &subnets,
			},
		},
	}
}

func connectPair(ctx *cli.Context) error {

	nameList := strings.SplitN(netns, ",", 2)
	subnetList := strings.SplitN(subnets, ",", 2)
	vethPeers, err := network.CraeteVethPair("", "")
	if err != nil {
		return err
	}
	// netns1
	if err := network.AssignIpAndUp(nameList[0], subnetList[0], vethPeers[0]); err != nil {
		return err
	}
	// netns2
	if err := network.AssignIpAndUp(nameList[1], subnetList[1], vethPeers[1]); err != nil {
		return err
	}

	data := struct {
		Name1   string
		Name2   string
		Subnet1 string
		Subnet2 string
		Veth1   string
		Veth2   string
	}{
		Name1:   nameList[0],
		Subnet1: subnetList[0],
		Veth1:   vethPeers[0],
		Name2:   nameList[1],
		Subnet2: subnetList[1],
		Veth2:   vethPeers[1],
	}

	connectPairTemplate.Execute(os.Stdout, &data)

	return nil
}
