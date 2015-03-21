package main

import (
	"os"

	"github.com/hashtock/tracker/cli"
)

func main() {
	app := cli.CliApp()
	app.Run(os.Args)
}
