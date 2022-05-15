package cgroups

import (
	"github.com/urfave/cli/v2"
	"log"
	"simple-container/pkg/cgroups/subsystems"
)

func NewDeleteCgroupCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-cgroup",
		Usage:     "delete cgroup",
		UsageText: "scadm [global options] delete-cgroup [options]",
		Action:    deleteCgroup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of cgroup",
				Destination: &name,
			},
		},
	}
}

func deleteCgroup(ctx *cli.Context) error {

	if err := subsystems.DeleteCgroup(name); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
