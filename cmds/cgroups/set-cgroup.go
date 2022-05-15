package cgroups

import (
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/pkg/cgroups"
	"strings"
	"text/template"
)

var setCgroupTemplate = template.Must(template.New("simple-container controller setCgroup").Parse(`
------------------------------------------------
simple-container controller:
    cgroup path: {{.Path}}
    cgroup limits: {{.Limits}}
------------------------------------------------
`))

var limits string

func NewSetCgroupCommand() *cli.Command {
	return &cli.Command{
		Name:      "set-cgroup",
		Usage:     "set cgroup limits",
		UsageText: "scadm [global options] create-cgroup [options]",
		Action:    setCgroup,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of cgroup",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "limits",
				Usage:       "cgroup limit name and value",
				Destination: &limits,
			},
		},
	}
}

func setCgroup(ctx *cli.Context) error {

	// pre-check
	if _, err := cgroups.FindCgroupPath(name); err != nil {
		cgroups.CreateCgroup(name)
	}

	splits := strings.Split(limits, ",")
	// txn
	for i := 0; i < len(splits); i++ {
		kv := strings.SplitN(splits[i], "=", 2)
		cName := strings.SplitN(name, ":", 2)[1]
		err := cgroups.SetCgroupLimit(kv[0], kv[1], cName)
		if err != nil {
			return err
		}
	}

	path, err := cgroups.FindCgroupPath(name)
	if err != nil {
		return err
	}
	data := struct {
		Path   string
		Limits []string
	}{
		Path:   path,
		Limits: splits,
	}

	setCgroupTemplate.Execute(os.Stdout, &data)

	return nil
}
