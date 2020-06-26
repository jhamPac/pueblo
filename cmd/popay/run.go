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
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ := cmd.Flags().GetUint64(flagPort)

			fmt.Println("Launching Popay node and its HTTP API...")

			bootstrap := node.NewPeerNode(
				"34.83.161.134",
				9000,
				true,
				false,
			)

			n := node.New(getDataDirFromCmd(cmd), ip, port, bootstrap)
			err := n.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	addDefaultRequiredFlags(cmd)
	cmd.Flags().String(flagIP, node.DefaultIP, "exposed IP for communication with peers")
	cmd.Flags().Uint64(flagPort, node.DefaultHTTPort, "exposed HTTP port for communication with peers")
	return cmd
}
