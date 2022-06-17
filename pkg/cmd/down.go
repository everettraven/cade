package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [WORKSPACE]",
	Short: "removes a containerized development workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return down(args[0])
	},
}

func down(workspaceName string) error {
	docker := containerutil.NewDockerUtil()

	container := containerutil.Container{
		Name: workspaceName,
	}

	fmt.Println("Stopping the workspace container:", workspaceName)
	out, err := docker.StopContainer(container)
	if err != nil {
		return fmt.Errorf("encountered an error stopping the workspace container: %w | out: %s", err, out)
	}

	fmt.Println("Removing the workspace container:", workspaceName)
	out, err = docker.RemoveContainer(container)
	if err != nil {
		return fmt.Errorf("encountered an error removing the workspace container: %w | out: %s", err, out)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("encountered an error getting the user home directory: %w", err)
	}

	workspaceDir := filepath.Join(home, ".cade", "tmp", workspaceName)
	fmt.Println("Cleaning up the workspace tmp directory:", workspaceDir)
	err = os.RemoveAll(workspaceDir)
	if err != nil {
		return fmt.Errorf("encountered an error removing the workspace tmp directory: %w", err)
	}

	return nil
}
