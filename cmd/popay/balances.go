package main

import (
	"fmt"
	"os"

	"github.com/jhampac/pueblo/database"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balances",
		Short: "interact with balances (list...)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(balancesListCmd())

	return cmd
}

func balancesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "lists all balances",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Accounts balances at %x:\n", state.LatestBlockHash())
			fmt.Println("------------------------")

			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account, balance))
			}
		},
	}

	addDefaultRequiredFlags(cmd)

	return cmd
}
