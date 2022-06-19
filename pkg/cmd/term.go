package cmd

import (
	"fmt"

	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
)

var termCmd = &cobra.Command{
	Use:   "term [WORKSPACE]",
	Short: "starts a terminal in the workspace specified",
	RunE: func(cmd *cobra.Command, args []string) error {
		return term(args[0], containerutil.NewContainerUtil())
	},
}

func term(workspaceName string, containerUtil containerutil.ContainerUtil) error {
	execOpts := containerutil.ExecOptions{
		Interactive: true,
		Tty:         true,
	}

	containerName := fmt.Sprintf("cade-workspace-%s", workspaceName)

	err := containerUtil.Exec(execOpts, containerName, "/bin/sh")
	if err != nil {
		return fmt.Errorf("encountered an error starting the workspace terminal: %w", err)
	}

	return nil
}
