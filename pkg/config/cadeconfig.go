package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type WorkspaceConfig struct {
	Prebuilt      string `json:"prebuilt" yaml:"prebuilt"`
	Containerfile string `json:"containerfile" yaml:"containerfile"`
	Workdir       string `json:"workdir" yaml:"workdir"`
	WorkspaceName string `json:"workspace_name" yaml:"workspace_name"`
	Context       string `json:"context" yaml:"context"`
}

// ParseWorkspaceConfig will parse a WorkspaceConfig from the provided source.
// The path can either be a URL or a local filepath
func ParseWorkspaceConfig(path string) (*WorkspaceConfig, error) {
	config := &WorkspaceConfig{}

	configBytes := []byte{}
	var err error

	if strings.Contains(path, "https://") {
		configBytes, err = fetchWorkspaceConfigFromURL(path)
	} else {
		configBytes, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("encountered an error reading the cade config: %w", err)
	}

	if filepath.Ext(path) == ".json" {
		json.Unmarshal(configBytes, &config)
	} else if filepath.Ext(path) == ".yaml" {
		yaml.Unmarshal(configBytes, config)
	} else {
		return nil, fmt.Errorf("unsupported config file type. must be one of JSON or YAML")
	}

	return config, nil
}

func fetchWorkspaceConfigFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("encountered an error fetching the data: %w", err)
	}

	bytes := make([]byte, response.ContentLength)
	_, err = response.Body.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("encountered an error reading the data: %w", err)
	}

	return bytes, nil
}
