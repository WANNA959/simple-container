package network

import (
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/pkg/network"
	"text/template"
)

var createNetnsTemplate = template.Must(template.New("simple-container controller createNetns").Parse(`
------------------------------------------------
simple-container controller:
    netns name: {{.Name}}
------------------------------------------------
`))

var name string
var subnet string

func NewCreateNetnsCommand() *cli.Command {
	return &cli.Command{
		Name:      "create-netns",
		Usage:     "create network namespace",
		UsageText: "scadm [global options] create-netns [options]",
		Action:    createNetns,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of netns",
				Destination: &name,
			},
		},
	}
}

func createNetns(ctx *cli.Context) error {

	if err := network.AddNetns(name); err != nil {
		return err
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	createNetnsTemplate.Execute(os.Stdout, &data)

	return nil
}
