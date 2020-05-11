package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Start() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Failed running %q\n", os.Args[1])
		os.Exit(1)
	}
}

var mainCmd = &cobra.Command{
	Use:          "generator-mobile",
	Short:        "generator-mobile",
	SilenceUsage: true,
}

func init() {
	cobra.EnableCommandSorting = false
	mainCmd.AddCommand(
		startCmd,
	)
}

func Run(args []string) error {
	mainCmd.SetArgs(args)
	return mainCmd.Execute()
}
