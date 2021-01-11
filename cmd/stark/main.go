package main

import (
	"log"
	"os"

	"github.com/Allenxuxu/stark/cmd/stark/service"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "stark",
		Usage:   "stark ctl",
		Version: "v0.0.1",
		Authors: []*cli.Author{{Name: "徐旭", Email: "120582243@qq.com"}},
	}

	app.Commands = append(app.Commands, service.Command())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
