package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "popay",
		Short: "Random name for a blockchain",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Hello from Popay") },
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(balancesCmd())
	rootCmd.AddCommand(txCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
