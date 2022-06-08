package container

import (
	"github.com/urfave/cli/v2"
	"log"
	"simple-container/pkg/cgroups/subsystems"
	"simple-container/pkg/container"
	"strings"
)

var tty bool
var containerName string
var limits string
var net string
var volume string

func NewRunCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "run a container",
		UsageText: "scadm [global options] run [options]",
		Action:    runWithCommand,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "it",
				Usage:       "enable tty",
				Destination: &tty,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "container name",
				Destination: &containerName,
			},
			&cli.StringFlag{
				Name:        "v",
				Usage:       "container volumn",
				Destination: &volume,
			},
			&cli.StringFlag{
				Name:        "limits",
				Usage:       "cpu and memory limit",
				Destination: &limits,
			},
			&cli.StringFlag{
				Name:        "net",
				Usage:       "network namespace name",
				Destination: &net,
			},
		},
	}
}

func runWithCommand(ctx *cli.Context) error {
	var cmdArray []string
	for i := 0; i < ctx.Args().Len(); i++ {
		cmdArray = append(cmdArray, ctx.Args().Get(i))
	}
	imageName := cmdArray[0]
	cpuLimitsMap := make(map[string]string)
	memoryLimitsMap := make(map[string]string)
	splits := strings.Split(limits, ",")
	for i := 0; i < len(splits); i++ {
		kv := strings.SplitN(splits[i], "=", 2)
		if strings.Contains(kv[0], "cpu") {
			cpuLimitsMap[kv[0]] = kv[1]
		} else if strings.Contains(kv[0], "memory") {
			memoryLimitsMap[kv[0]] = kv[1]
		}
	}
	resConf := &subsystems.ResourceConfig{
		CpuLimits:    cpuLimitsMap,
		MemoryLimits: memoryLimitsMap,
	}
	if err := container.RunWithCommand(tty, resConf, volume, containerName, imageName, net); err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

func run(ctx *cli.Context) error {

	var cmdArray []string
	args := ctx.Args()
	for i := 0; i < args.Len(); i++ {
		cmdArray = append(cmdArray, args.Get(i))
	}
	log.Printf("cmdArray:%+v", cmdArray)
	imageName := cmdArray[0]
	cpuLimitsMap := make(map[string]string)
	memoryLimitsMap := make(map[string]string)
	splits := strings.Split(limits, ",")
	for i := 0; i < len(splits); i++ {
		kv := strings.SplitN(splits[i], "=", 2)
		if strings.Contains(kv[0], "cpu") {
			cpuLimitsMap[kv[0]] = kv[1]
		} else if strings.Contains(kv[0], "memory") {
			memoryLimitsMap[kv[0]] = kv[1]
		}
	}
	resConf := &subsystems.ResourceConfig{
		CpuLimits:    cpuLimitsMap,
		MemoryLimits: memoryLimitsMap,
	}
	if err := container.Run(tty, cmdArray, resConf, volume, containerName, imageName, net); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
