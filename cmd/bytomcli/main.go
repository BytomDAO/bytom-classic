package main

import (
	"runtime"

	cmd "github.com/bytom/bytom-classic/cmd/bytomcli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
