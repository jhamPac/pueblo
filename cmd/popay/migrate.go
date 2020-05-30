package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jhampac/pueblo/database"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate tx to block database",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		state.Close()

		block0 := database.NewBlock(
			database.Hash{},
			uint64(time.Now().Unix()),
			[]database.Tx{
				database.NewTx("redcloud", "redcloud", 3, ""),
				database.NewTx("redcloud", "redcloud", 700, "reward"),
			},
		)

		state.AddBlock(block0)
		block0hash, _ := state.Persist()

		block1 := database.NewBlock(
			block0hash,
			uint64(time.Now().Unix()),
			[]database.Tx{
				database.NewTx("redcloud", "sittingbull", 2000, ""),
				database.NewTx("redcloud", "redcloud", 100, "reward"),
				database.NewTx("sittingbull", "redcloud", 10, ""),
				database.NewTx("sittingbull", "woundedknee", 1000, ""),
				database.NewTx("sittingbull", "redcloud", 50, ""),
				database.NewTx("redcloud", "redcloud", 700, "reward"),
			},
		)

		state.AddBlock(block1)

		_, err = state.Persist()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
