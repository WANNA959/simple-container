package network

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"simple-container/pkg/network"
	"text/template"
)

var deleteNetnsTemplate = template.Must(template.New("simple-container controller deleteNetns").Parse(`
------------------------------------------------
simple-container controller:
    netns name: {{.Name}}
------------------------------------------------
`))

func NewDeleteNetnsCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete-netns",
		Usage:     "delete network namespace",
		UsageText: "scadm [global options] create-netns [options]",
		Action:    deleteNetns,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of netns",
				Destination: &name,
			},
		},
	}
}

func deleteNetns(ctx *cli.Context) error {

	if err := network.DeleteNetns(name); err != nil {
		log.Fatalln(err)
		return err
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	deleteNetnsTemplate.Execute(os.Stdout, &data)

	return nil
}
