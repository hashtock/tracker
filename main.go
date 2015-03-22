package main

import (
	"os"

	"github.com/hashtock/tracker/cli"
)

func main() {
	app := cli.App()
	app.Run(os.Args)
}
