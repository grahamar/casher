package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/grahamar/casher/root"

	// commands
	_ "github.com/grahamar/casher/add"
	_ "github.com/grahamar/casher/fetch"
	_ "github.com/grahamar/casher/push"
)

func main() {
	log.SetHandler(cli.Default)

	args := os.Args[1:]
	root.Command.SetArgs(args)

	if err := root.Command.Execute(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
