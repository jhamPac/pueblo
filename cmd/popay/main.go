package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const flagDataDir = "datadir"
const flagPort = "port"

func main() {
	rootCmd := &cobra.Command{
		Use:   "popay",
		Short: "Popay your friends and family fast!",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Popay your friends!") },
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(txCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(runCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// so we don't have to repeat this code over and over in other cmd files
// just call this function and pass in a pointer to said cmd you want to attach dataDir command to
func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
