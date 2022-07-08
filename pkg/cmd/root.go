package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "cade",
	Short: "cade is a CLI tool for using Containers as Development Environments",
	Example: `
	## Starting a workspace 
	cade up https://raw.githubusercontent.com/everettraven/cade/main/example/cadeconfig.yaml

	## Starting a terminal in a workspace
	cade term cade-test

	## Stopping a workspace
	cade down cade-test

	## Get the current cade version
	cade version
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(termCmd)
	rootCmd.AddCommand(listCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
