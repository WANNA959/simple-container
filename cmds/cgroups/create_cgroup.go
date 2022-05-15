package cgroups

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"simple-container/pkg/cgroups/subsystems"
	"text/template"
)

var createCgroupTemplate = template.Must(template.New("simple-container controller createCgroup").Parse(`
------------------------------------------------
simple-container controller:
    cgroup path: {{.Path}}
------------------------------------------------
`))

var name string

func NewCreateCgroupCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-cgroup",
		Usage:     "create cgroup",
		UsageText: "scadm [global options] create-cgroup [options]",
		Action:    createCgroup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of cgroup",
				Destination: &name,
			},
		},
	}
}

func createCgroup(ctx *cli.Context) error {

	if err := subsystems.CreateCgroup(name); err != nil {
		log.Fatalln(err)
		return err
	}

	path, err := subsystems.FindCgroupPath(name)
	if err != nil {
		return err
	}
	data := struct {
		Path string
	}{
		Path: path,
	}

	createCgroupTemplate.Execute(os.Stdout, &data)

	return nil
}
