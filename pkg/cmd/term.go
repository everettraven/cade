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
		return term(args[0])
	},
}

// TODO(everettraven): make this a bit more configurable, maybe allow for
// choosing the shell to start up with.
func term(workspaceName string) error {
	docker := containerutil.NewDockerUtil()

	execOpts := containerutil.ExecOptions{
		Interactive: true,
		Tty:         true,
	}

	err := docker.Exec(execOpts, workspaceName, "/bin/bash")
	if err != nil {
		return fmt.Errorf("encountered an error starting the workspace terminal: %w", err)
	}

	return nil
}
