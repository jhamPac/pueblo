package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "pueblo",
		Short: "Random name for a blockchain",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Hello from Popay") },
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
