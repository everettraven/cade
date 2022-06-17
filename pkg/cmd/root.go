package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "cade",
	Short: "cade is a CLI tool for using Containers as Development Environments",
	Long:  "cade is a CLI tool that easily enables the use of Containers as Development Environments on a local environment.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(termCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
