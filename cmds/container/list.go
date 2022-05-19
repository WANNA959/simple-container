package container

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"simple-container/pkg/container"
	"simple-container/pkg/utils"
	"text/tabwriter"
)

func NeListCommand() *cli.Command {
	return &cli.Command{
		Name:      "ps",
		Usage:     "connect to host bridge",
		UsageText: "scadm [global options] ps",
		Action:    listContainer,
	}
}

func listContainer(ctx *cli.Context) error {

	pids := utils.GetAllPid()
	var containers []*container.ContainerInfo
	for _, pid := range pids {
		info := container.GetContainerInfo(pid)
		containers = append(containers, info)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for idx, item := range containers {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			idx+1,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime,
		)
	}
	if err := w.Flush(); err != nil {
		log.Fatalf("Flush error %v", err)
		return err
	}
	return nil
}
