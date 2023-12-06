package main

import (
	"log"
	"os"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/go-pandoc/config"
	_ "github.com/pschlump/go-pandoc/pandoc/fetcher/data"
	_ "github.com/pschlump/go-pandoc/pandoc/fetcher/http"
	_ "github.com/pschlump/go-pandoc/pandoc/fetcher/redis"
	"github.com/pschlump/go-pandoc/server"
	"github.com/urfave/cli/v2"
)

func main() {

	var err error
	defer func() {
		if err != nil {
			log.Printf("[go-pandoc]: %s\n", err.Error())
		}
	}()

	app := cli.NewApp()

	app.Usage = "A server for pandoc command"
	app.Commands = cli.Commands{
		&cli.Command{
			Name:   "run",
			Usage:  "run pandoc service",
			Action: run,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "config filename",
					Value:   "app.conf",
				},
				&cli.StringFlag{
					Name:  "cwd",
					Usage: "change work dir",
				},
			},
		},
	}

	err = app.Run(os.Args)
}

func run(ctx *cli.Context) (err error) {

	cwd := ctx.String("cwd")
	if len(cwd) != 0 {
		err = os.Chdir(cwd)
	}

	if err != nil {
		return
	}

	configFile := ctx.String("config")

	conf := config.NewConfig(
		config.ConfigFile(configFile),
	)

	authkey := conf.GetString("auth-key", "ya, no api key specified")
	dbgo.Printf("Auth Key Is: %(cyan)->%s<-\n", authkey)

	srv, err := server.New(conf)

	if err != nil {
		return
	}

	err = srv.Run()

	return
}
