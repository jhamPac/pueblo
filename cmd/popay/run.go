package main

import (
	"fmt"
	"os"

	"github.com/jhampac/pueblo/node"
	"github.com/spf13/cobra"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launches the Popay node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			port, _ := cmd.Flags().GetUint64(flagPort)

			fmt.Println("Launching Popay node and its HTTP API...")

			bootstrap := node.NewPeerNode(
				"34.83.161.134",
				9000,
				true,
				true,
			)

			n := node.New(dataDir, port, bootstrap)

			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(cmd)
	cmd.Flags().Uint64(flagPort, node.DefaultHTTPort, "exposed HTTP port for communication with peers")
	return cmd
}
