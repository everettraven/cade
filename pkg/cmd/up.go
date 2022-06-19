package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/cade/pkg/config"
	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
)

var name string
var build bool

var upCmd = &cobra.Command{
	Use:   "up [WORKSPACE]",
	Short: "creates a containerized development workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		return up(args[0])
	},
}

func init() {
	upCmd.Flags().StringVarP(&name, "name", "n", "", "sets the workspace name")
	upCmd.Flags().BoolVarP(&build, "build", "b", false, "force the workspace image to be built")
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

	if workspaceConfig.Prebuilt == "" || build {
		context := "."
		if workspaceConfig.Context != "" {
			context = workspaceConfig.Context
		}
		fmt.Println("Building the image (this could take some time...)")
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

	baseWorkspaceDir := filepath.Join(home, "cade", "workspaces")
	fmt.Println("Ensuring the", baseWorkspaceDir, "directory is created")
	err = os.MkdirAll(baseWorkspaceDir, 0777)
	if err != nil {
		return fmt.Errorf("encountered an error ensuring the directory `%s` exists: %w", baseWorkspaceDir, err)
	}

	workspaceDir := filepath.Join(baseWorkspaceDir, wkspName)

	volumes := []containerutil.Volume{
		{
			HostPath:  workspaceDir,
			MountPath: workspaceConfig.Workdir,
		},
	}

	if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
		fmt.Println("Copying files from container to workspace directory")
		out, err := docker.CopyToHost(container, volumes[0])
		if err != nil {
			return fmt.Errorf("encountered an error copying files from container to host: %w | out: %s", err, out)
		}
	} else if err != nil {
		return fmt.Errorf("encountered an error checking if directory `%s` already exists: %w", workspaceDir, err)
	}

	fmt.Println("Running the workspace container")
	out, err := docker.Run(container, volumes)
	if err != nil {
		return fmt.Errorf("encountered an error running the workspace image: %w | out: %s", err, out)
	}

	fmt.Println("Workspace ready! The workspace name is", wkspName, "and the mounted working directory is", workspaceDir)
	return nil
}
