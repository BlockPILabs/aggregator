package commands

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
)

func RootApp() *cli.App {
	runCmd := RunCommand()
	cmd := &cli.App{
		Name:                 "aggregator",
		Usage:                "RPCHub aggregator",
		Flags:                runCmd.Flags,
		EnableBashCompletion: true,
		BashComplete:         cli.DefaultAppComplete,
		Before: func(cli *cli.Context) error {
			if err := initConfig(cli); err != nil {
				return err
			}
			return nil
		},
		Action: func(cli *cli.Context) error {
			return runCommand(cli, "run", os.Args[1:]...)
		},
		Commands: []*cli.Command{
			runCmd,
			InitCommand(),
		},
	}
	return cmd
}

func runCommand(app *cli.Context, cmd string, args ...string) error {
	args = append([]string{app.App.Name, cmd}, args...)

	subCommand := app.App.Command(cmd)
	if subCommand == nil {
		return errors.New("not sub command : " + cmd)
	}

	subCommand.SkipFlagParsing = true
	subCommand.Flags = []cli.Flag{}
	return subCommand.Run(app)
}

func initConfig(cli *cli.Context) error {
	//viper.SetConfigFile(cli.String(config.FlagConfigFile.Name))

	//return config.LoadConfig()
	return nil
}
