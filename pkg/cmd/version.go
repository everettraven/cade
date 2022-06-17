package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "output the version of the currently installed cade",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
