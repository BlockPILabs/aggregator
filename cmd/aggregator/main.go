package main

import (
	"errors"
	"fmt"
	"github.com/BlockPILabs/aggregator/cmd/aggregator/commands"
	"github.com/BlockPILabs/aggregator/version"
	"os"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	println(version.Version)

	time.Local = time.UTC
	initPath()

	app := commands.RootApp()
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}

func initPath() {
	dir, err := gfile.Home(".blockpi/aggregator")
	if err != nil {
		panic(err)
	}

	if !gfile.Exists(dir) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			panic(err)
		}
	}

	if !gfile.IsDir(dir) {
		panic(errors.New(fmt.Sprintf("%s is not a dir", dir)))
	}

	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
