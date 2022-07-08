package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/everettraven/cade/pkg/containerutil"
	yaml "gopkg.in/yaml.v2"
)

type WorkspaceConfig struct {
	Prebuilt      string                 `json:"prebuilt" yaml:"prebuilt"`
	Containerfile string                 `json:"containerfile" yaml:"containerfile"`
	Workdir       string                 `json:"workdir" yaml:"workdir"`
	WorkspaceName string                 `json:"workspace_name" yaml:"workspace_name"`
	Context       string                 `json:"context" yaml:"context"`
	Volumes       []containerutil.Volume `json: "volumes" yaml:"volumes"`
	Network       string                 `json:"network" yaml:"network"`
}

// ParseWorkspaceConfig will parse a WorkspaceConfig from the provided source.
// The path can either be a URL or a local filepath
func ParseWorkspaceConfig(path string) (*WorkspaceConfig, error) {
	config := &WorkspaceConfig{}

	var configBytes []byte
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

	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("encountered an error reading the data: %w", err)
	}

	return bytes, nil
}
