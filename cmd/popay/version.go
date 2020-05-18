package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const major = "0"
const minor = "1"
const patch = "0"
const description = "TX Add && Balances List"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print popay version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s.%s.%s-beta %s", major, minor, patch, description)
	},
}
