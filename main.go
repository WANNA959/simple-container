package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"simple-container/cmds"
)

func main() {
	app := cmds.NewApp()
	app.Commands = []*cli.Command{}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error options: %s\n", err.Error())
		os.Exit(-1)
	}
}
