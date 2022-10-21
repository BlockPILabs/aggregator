package main

import (
	"aggregator/cmd/aggregator/commands"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
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
		err = os.MkdirAll(dir, 0600)
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
