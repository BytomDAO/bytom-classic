package commands

import (
	"runtime"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/anonimitycash/anonimitycash-classic/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Anonimitycashcli",
	Run: func(cmd *cobra.Command, args []string) {
		jww.FEEDBACK.Printf("Anonimitycashcli v%s %s/%s\n", version.Version, runtime.GOOS, runtime.GOARCH)
	},
}
