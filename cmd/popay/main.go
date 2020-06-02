package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const flagDataDir = "datadir"

func main() {
	rootCmd := &cobra.Command{
		Use:   "popay",
		Short: "Random name for a blockchain",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Popay your friends!") },
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(txCmd())
	rootCmd.AddCommand(migrateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute pat to the node data dir where the DB is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
