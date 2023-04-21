package commands

import (
	"github.com/BlockPILabs/aggregator/config"
	"github.com/BlockPILabs/aggregator/loadbalance"
	"github.com/BlockPILabs/aggregator/middleware"
	"github.com/BlockPILabs/aggregator/middleware/plugins"
	"github.com/BlockPILabs/aggregator/server"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func RunCommand() *cli.Command {
	return &cli.Command{
		Name:    "run",
		Aliases: []string{"start"},
		Flags:   append([]cli.Flag{}, InitCommand().Flags...),
		Before: func(cli *cli.Context) error {
			err := runCommand(cli, "init")
			if err != nil {
				return err
			}

			config.Load()

			loadbalance.LoadFromConfig()

			middleware.Append(
				plugins.NewRequestValidatorMiddleware(),
				plugins.NewSafetyMiddleware(),
				plugins.NewLoadBalanceMiddleware(),
				plugins.NewHttpProxyMiddleware(),
				plugins.NewCorsMiddleware(),
			)

			return nil
		},
		Action: func(context *cli.Context) error {
			wg := errgroup.Group{}
			wg.Go(func() error {
				return server.NewManageServer()
			})
			wg.Go(func() error {
				return server.NewServer()
			})
			return wg.Wait()
		},
		Subcommands: []*cli.Command{},
	}

}
