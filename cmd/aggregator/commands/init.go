package commands

import (
	"github.com/urfave/cli/v2"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Flags: []cli.Flag{},
		Before: func(context *cli.Context) error {
			return nil
		},
		Action: func(cli *cli.Context) error {

			return nil
		},
	}
}
