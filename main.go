package main

import (
	"os"

	"jdiff/internal/cli"
)

func main() {
	app := cli.New(os.Stdout, os.Stderr)
	exitCode := app.Run(os.Args[1:])
	os.Exit(exitCode)
}
