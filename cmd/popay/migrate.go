package main

import (
	"github.com/spf13/cobra"
	"github.com/jhampac/pueblo/database"
)

var migrate = &cobra.Command{
	Use: "migrate",
	Short: "migrate tx to block database",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}