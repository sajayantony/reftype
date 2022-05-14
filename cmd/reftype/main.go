package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:          "reftype [command]",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		refsCmd(),
		manifestCmd(),
		pushCmd(),
	)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
