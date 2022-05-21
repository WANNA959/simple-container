package container

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"simple-container/pkg/container"
	"text/tabwriter"
)

func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:      "ps",
		Usage:     "connect to host bridge",
		UsageText: "scadm [global options] ps",
		Action:    listContainer,
	}
}

func listContainer(ctx *cli.Context) error {

	dirURL := container.DefaultInfoLocation
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		log.Fatalf("Read dir %s error %v", dirURL, err)
		return err
	}

	var containers []*container.ContainerInfo
	for _, file := range files {
		tmpContainer, err := container.GetContainerInfo(file)
		if err != nil {
			log.Fatalf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
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
