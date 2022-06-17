package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/cade/pkg/config"
	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var name string

var upCmd = &cobra.Command{
	Use:   "up [OPTIONS] [WORKSPACE]",
	Short: "creates a containerized development workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return up(args[0])
	},
}

func init() {
	pflag.StringVarP(&name, "name", "n", "", "sets the workspace name")
}

func up(workspace string) error {
	fmt.Println("Parsing the workspace configuration file")
	workspaceConfig, err := config.ParseWorkspaceConfig(workspace)
	if err != nil {
		return fmt.Errorf("encountered an error getting the cade config: %w", err)
	}

	wkspName := workspaceConfig.WorkspaceName

	if name != "" {
		wkspName = name
	}

	fmt.Println("Creating containerized workspace:", wkspName)

	docker := containerutil.NewDockerUtil()

	if workspaceConfig.Prebuilt == "" {
		context := "."
		if workspaceConfig.Context != "" {
			context = workspaceConfig.Context
		}
		fmt.Println("No prebuilt image found, building the image")
		out, err := docker.Build(workspaceConfig.Containerfile, wkspName, context)
		if err != nil {
			return fmt.Errorf("encountered an error building the workspace image: %w | out: %s", err, out)
		}

		workspaceConfig.Prebuilt = wkspName
	}

	container := containerutil.Container{
		Name:  wkspName,
		Image: workspaceConfig.Prebuilt,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("encountered an error getting the user home directory: %w", err)
	}

	workspaceDir := filepath.Join(home, ".cade", "tmp", wkspName)

	fmt.Println("Creating temporary directory to act as working directory:", workspaceDir)
	err = os.MkdirAll(workspaceDir, 0777)
	if err != nil {
		return fmt.Errorf("encountered an error creating volume directory: %w", err)
	}

	volumes := []containerutil.Volume{
		{
			HostPath:  workspaceDir,
			MountPath: workspaceConfig.Workdir,
		},
	}

	fmt.Println("Running the workspace container")
	out, err := docker.Run(container, volumes)
	if err != nil {
		return fmt.Errorf("encountered an error running the workspace image: %w | out: %s", err, out)
	}

	fmt.Println("Workspace ready! The workspace name is", wkspName, "and the mounted working directory is", workspaceDir)
	return nil
}
