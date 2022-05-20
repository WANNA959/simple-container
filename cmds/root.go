package cmds

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
	"simple-container/pkg/container"
	"simple-container/pkg/network"
	"simple-container/pkg/sqlite"
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
	app.Before = initBridge
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s\n", app.Name, app.Version)
		fmt.Printf("go version %s\n", runtime.Version())
	}

	return app
}

var initBridge = func(ctx *cli.Context) error {
	sqlite.InitSqlite()
	os.MkdirAll(container.DefaultInfoLocation, 622)
	if err := network.GenerateBridgeOrSkip(network.DefaultBridgeSubnet, network.DefaultDocker0Bridge); err != nil {
		return err
	}
	return nil
}
