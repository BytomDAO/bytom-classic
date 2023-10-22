package main

import (
	"runtime"

	cmd "github.com/anonimitycash/anonimitycash-classic/cmd/anonimitycashcli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
