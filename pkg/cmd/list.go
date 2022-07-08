package cmd

import (
	"fmt"
	"strings"

	"github.com/everettraven/cade/pkg/containerutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the current workspaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		return list(containerutil.NewContainerUtil())
	},
}

func list(containerUtil containerutil.ContainerUtil) error {
	containers, err := containerUtil.ContainerList()
	if err != nil {
		return fmt.Errorf("encountered an error attempting to get a list of containers: %w", err)
	}

	fmt.Println("Available workspaces:")
	for _, container := range containers {
		if strings.Contains(container.Name, "cade-workspace") {
			fmt.Println("-", strings.Replace(container.Name, "cade-workspace-", "", -1))
		}
	}

	return nil
}
