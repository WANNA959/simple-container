package cmds

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	"simple-container/pkg/version"
)

type AccessConfig struct {
}

var GlobalConfig AccessConfig

var homeDir string = func() string {
	if home, err := os.UserHomeDir(); err != nil {
		return ""
	} else {
		return home
	}
}()

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "scadm"
	app.Usage = "scadm, a commond-line tool to control simple container"
	app.Version = version.Version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s\n", app.Name, app.Version)
		fmt.Printf("go version %s\n", runtime.Version())
	}

	return app
}
