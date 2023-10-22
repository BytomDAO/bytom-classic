package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/anonimitycash/anonimitycash-classic/util"
)

// anonimitycashcli usage template
var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletDisable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

  available with wallet enable:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletEnable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func newSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func newSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}

// AnonimitycashcliCmd is Anonimitycashcli's root command.
// Every other command attached to AnonimitycashcliCmd is a child command to it.
var AnonimitycashcliCmd = &cobra.Command{
	Use:   "anonimitycashcli",
	Short: "Anonimitycashcli is a commond line client for anonimitycash core (a.k.a. anonimitycashd)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.SetUsageTemplate(usageTemplate)
			cmd.Usage()
		}
	},
}

// Execute adds all child commands to the root command AnonimitycashcliCmd and sets flags appropriately.
func Execute() {

	AddCommands()
	AddTemplateFunc()

	if _, err := AnonimitycashcliCmd.ExecuteC(); err != nil {
		os.Exit(util.ErrLocalExe)
	}
}

// AddCommands adds child commands to the root command AnonimitycashcliCmd.
func AddCommands() {
	AnonimitycashcliCmd.AddCommand(createAccessTokenCmd)
	AnonimitycashcliCmd.AddCommand(listAccessTokenCmd)
	AnonimitycashcliCmd.AddCommand(deleteAccessTokenCmd)
	AnonimitycashcliCmd.AddCommand(checkAccessTokenCmd)

	AnonimitycashcliCmd.AddCommand(createAccountCmd)
	AnonimitycashcliCmd.AddCommand(deleteAccountCmd)
	AnonimitycashcliCmd.AddCommand(listAccountsCmd)
	AnonimitycashcliCmd.AddCommand(updateAccountAliasCmd)
	AnonimitycashcliCmd.AddCommand(createAccountReceiverCmd)
	AnonimitycashcliCmd.AddCommand(listAddressesCmd)
	AnonimitycashcliCmd.AddCommand(validateAddressCmd)
	AnonimitycashcliCmd.AddCommand(listPubKeysCmd)

	AnonimitycashcliCmd.AddCommand(createAssetCmd)
	AnonimitycashcliCmd.AddCommand(getAssetCmd)
	AnonimitycashcliCmd.AddCommand(listAssetsCmd)
	AnonimitycashcliCmd.AddCommand(updateAssetAliasCmd)

	AnonimitycashcliCmd.AddCommand(getTransactionCmd)
	AnonimitycashcliCmd.AddCommand(listTransactionsCmd)

	AnonimitycashcliCmd.AddCommand(getUnconfirmedTransactionCmd)
	AnonimitycashcliCmd.AddCommand(listUnconfirmedTransactionsCmd)
	AnonimitycashcliCmd.AddCommand(decodeRawTransactionCmd)

	AnonimitycashcliCmd.AddCommand(listUnspentOutputsCmd)
	AnonimitycashcliCmd.AddCommand(listBalancesCmd)

	AnonimitycashcliCmd.AddCommand(rescanWalletCmd)
	AnonimitycashcliCmd.AddCommand(walletInfoCmd)

	AnonimitycashcliCmd.AddCommand(buildTransactionCmd)
	AnonimitycashcliCmd.AddCommand(signTransactionCmd)
	AnonimitycashcliCmd.AddCommand(submitTransactionCmd)
	AnonimitycashcliCmd.AddCommand(estimateTransactionGasCmd)

	AnonimitycashcliCmd.AddCommand(getBlockCountCmd)
	AnonimitycashcliCmd.AddCommand(getBlockHashCmd)
	AnonimitycashcliCmd.AddCommand(getBlockCmd)
	AnonimitycashcliCmd.AddCommand(getBlockHeaderCmd)
	AnonimitycashcliCmd.AddCommand(getDifficultyCmd)
	AnonimitycashcliCmd.AddCommand(getHashRateCmd)

	AnonimitycashcliCmd.AddCommand(createKeyCmd)
	AnonimitycashcliCmd.AddCommand(deleteKeyCmd)
	AnonimitycashcliCmd.AddCommand(listKeysCmd)
	AnonimitycashcliCmd.AddCommand(updateKeyAliasCmd)
	AnonimitycashcliCmd.AddCommand(resetKeyPwdCmd)
	AnonimitycashcliCmd.AddCommand(checkKeyPwdCmd)

	AnonimitycashcliCmd.AddCommand(signMsgCmd)
	AnonimitycashcliCmd.AddCommand(verifyMsgCmd)
	AnonimitycashcliCmd.AddCommand(decodeProgCmd)

	AnonimitycashcliCmd.AddCommand(createTransactionFeedCmd)
	AnonimitycashcliCmd.AddCommand(listTransactionFeedsCmd)
	AnonimitycashcliCmd.AddCommand(deleteTransactionFeedCmd)
	AnonimitycashcliCmd.AddCommand(getTransactionFeedCmd)
	AnonimitycashcliCmd.AddCommand(updateTransactionFeedCmd)

	AnonimitycashcliCmd.AddCommand(isMiningCmd)
	AnonimitycashcliCmd.AddCommand(setMiningCmd)

	AnonimitycashcliCmd.AddCommand(netInfoCmd)
	AnonimitycashcliCmd.AddCommand(gasRateCmd)

	AnonimitycashcliCmd.AddCommand(versionCmd)
}

// AddTemplateFunc adds usage template to the root command AnonimitycashcliCmd.
func AddTemplateFunc() {
	walletEnableCmd := []string{
		createAccountCmd.Name(),
		listAccountsCmd.Name(),
		deleteAccountCmd.Name(),
		updateAccountAliasCmd.Name(),
		createAccountReceiverCmd.Name(),
		listAddressesCmd.Name(),
		validateAddressCmd.Name(),
		listPubKeysCmd.Name(),

		createAssetCmd.Name(),
		getAssetCmd.Name(),
		listAssetsCmd.Name(),
		updateAssetAliasCmd.Name(),

		createKeyCmd.Name(),
		deleteKeyCmd.Name(),
		listKeysCmd.Name(),
		resetKeyPwdCmd.Name(),
		checkKeyPwdCmd.Name(),
		signMsgCmd.Name(),

		buildTransactionCmd.Name(),
		signTransactionCmd.Name(),

		getTransactionCmd.Name(),
		listTransactionsCmd.Name(),
		listUnspentOutputsCmd.Name(),
		listBalancesCmd.Name(),

		rescanWalletCmd.Name(),
		walletInfoCmd.Name(),
	}

	cobra.AddTemplateFunc("WalletEnable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return true
			}
		}
		return false
	})

	cobra.AddTemplateFunc("WalletDisable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return false
			}
		}
		return true
	})
}
