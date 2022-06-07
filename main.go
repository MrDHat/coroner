package main

import (
	"coroner/cron"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func main() {
	clientApp := cli.NewApp()
	clientApp.Name = "coroner"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:  "parse",
			Usage: "Parse the cron expression",
			Action: func(c *cli.Context) error {
				args := c.Args()
				if len(args) == 0 {
					fmt.Println("Missing cron expression")
					return nil
				}
				cron.NewParser().Parse(args.Get(0))
				return nil
			},
		},
	}
	if err := clientApp.Run(os.Args); err != nil {
		panic(err)
	}
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(
		signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
}
