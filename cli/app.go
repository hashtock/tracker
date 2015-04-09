package cli

import (
	"github.com/codegangsta/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "tracker"
	app.Usage = "Twitter hashtag count tracking"
	app.Author = "Karol Dulęba"
	app.Email = "mr.fuxi@gmail.com"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose", Usage: "be more verbose"},
		cli.StringFlag{Name: "remote", Usage: "execute on remote server"},
	}
	app.Action = cmdWebAPI

	app.Commands = []cli.Command{
		{
			Name:      "web",
			ShortName: "w",
			Usage:     "run web server",
			Action:    cmdWebAPI,
		},
		{
			Name:      "listen",
			ShortName: "l",
			Usage:     "start counting tweets with requested hashtags",
			Action:    cmdListen,
		},
		{
			Name:      "tags",
			ShortName: "t",
			Usage:     "list current tags",
			Action:    cmdListTags,
		},
		{
			Name:      "counts",
			ShortName: "c",
			Usage:     "list counts - sum. Example: c 2h45m",
			Action:    cmdListTagCounts,
		},
		{
			Name:      "counts_detailed",
			ShortName: "cd",
			Usage:     "list counts - all data points. Example: cd 2h45m",
			Action:    cmdListTagCountsDetails,
		},
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "list current tags",
			Action:    cmdAddTags,
		},
		{
			Name:  "clear",
			Usage: "removes data (all|counts)",
			Subcommands: []cli.Command{
				{
					Name:   "all",
					Usage:  "removes all data",
					Action: cmdClearAll,
				},
				{
					Name:   "counts",
					Usage:  "removes all counts",
					Action: cmdClearCounts,
				},
			},
		},
		{
			Name:   "config",
			Usage:  "Shows configuration details",
			Action: cmdPrintConfHelp,
		},
	}
	return app
}
