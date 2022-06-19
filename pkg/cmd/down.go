package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
)

var persistWorkdir bool

var downCmd = &cobra.Command{
	Use:   "down [WORKSPACE]",
	Short: "removes a containerized development workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return down(args[0], containerutil.NewContainerUtil())
	},
}

func init() {
	downCmd.Flags().BoolVarP(&persistWorkdir, "persist-workdir", "p", false, "Persist the working directory of the workspace")
}

func down(workspaceName string, containerUtil containerutil.ContainerUtil) error {
	container := containerutil.Container{
		Name: fmt.Sprintf("cade-workspace-%s", workspaceName),
	}

	fmt.Println("Stopping the workspace container:", container.Name)
	out, err := containerUtil.StopContainer(container)
	if err != nil {
		return fmt.Errorf("encountered an error stopping the workspace container: %w | out: %s", err, out)
	}

	fmt.Println("Removing the workspace container:", container.Name)
	out, err = containerUtil.RemoveContainer(container)
	if err != nil {
		return fmt.Errorf("encountered an error removing the workspace container: %w | out: %s", err, out)
	}

	if !persistWorkdir {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("encountered an error getting the user home directory: %w", err)
		}

		workspaceDir := filepath.Join(home, "cade", "workspaces", workspaceName)
		fmt.Println("Cleaning up the workspace working directory:", workspaceDir)
		err = os.RemoveAll(workspaceDir)
		if err != nil {
			return fmt.Errorf("encountered an error removing the workspace tmp directory: %w", err)
		}
	}

	return nil
}
