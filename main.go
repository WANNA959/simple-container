package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/cmds"
	"simple-container/cmds/cgroups"
	"simple-container/cmds/container"
	"simple-container/cmds/network"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{
		// run container
		container.NewRunCommand(),
		container.NewListCommand(),
		// network namespace related
		network.NewCreateNetnsCommand(),
		network.NewConnectPairCommand(),
		network.NewConnectBridgeCommand(),
		network.NewDeleteNetnsCommand(),
		// cgroup related
		cgroups.NewCreateCgroupCommand(),
		cgroups.NewDeleteCgroupCommand(),
		cgroups.NewSetCgroupCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error options: %s\n", err.Error())
		os.Exit(-1)
	}
}
