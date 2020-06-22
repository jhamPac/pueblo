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
		cwd, _ := os.Getwd()
		state, err := database.NewStateFromDisk(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()

		block0 := database.NewBlock(
			database.Hash{},
			1,
			uint64(time.Now().Unix()),
			[]database.Tx{
				database.NewTx("alice", "bob", 100000, ""),
				database.NewTx("alice", "eve", 100000, ""),
			},
		)

		state.AddBlock(block0)
		block0hash, _ := state.Persist()

		block1 := database.NewBlock(
			block0hash,
			2,
			uint64(time.Now().Unix()),
			[]database.Tx{
				database.NewTx("alice", "alice", 1000, "reward"),
				database.NewTx("alice", "alice", 1000, "reward"),
				database.NewTx("alice", "alice", 1000, "reward"),
				database.NewTx("alice", "alice", 1000, "reward"),
				database.NewTx("alice", "alice", 1000, "reward"),
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
