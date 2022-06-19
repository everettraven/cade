package containerutil

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Docker struct{}

type dockerContainer struct {
	Id      string `json:"ID"`
	Name    string `json:"Names"`
	Image   string `json:"Image"`
	Command string `json:"Command"`
	Created string `json:"RunningFor"`
	Status  string `json:"Status"`
	State   string `json:"State"`
	Ports   string `json:"Ports"`
}

type dockerContainerList struct {
	Containers []dockerContainer
}

type dockerImage struct {
	Containers   string
	CreatedAt    string
	CreatedSince string
	Digest       string
	ID           string
	Repository   string
	SharedSize   string
	Size         string
	Tag          string
	UniqueSize   string
	VirtualSize  string
}

type dockerImageList struct {
	Images []dockerImage
}

// NewDockerUtil returns a ContainerUtil implementation
// that uses Docker as the container runtime
func NewDockerUtil() *Docker {
	return &Docker{}
}

// Run runs a container using the provided image and sets the container
// name to be the provided name. It also mounts any volumes provided.
// Returns an error if any occur during the process
func (d *Docker) Run(container Container, volumes []Volume, runArgs ...string) ([]byte, error) {
	args := []string{
		"run",
		"-d",
		"-t",
		"--name",
		container.Name,
	}

	for _, volume := range volumes {
		args = append(args, "-v", fmt.Sprintf("%s:%s", volume.HostPath, volume.MountPath))
	}

	args = append(args, container.Image)

	args = append(args, runArgs...)

	return runDockerCmd(args...)
}

// Build builds an image from the provided containerfile and tags it
// with the provided tag. Returns an error if any occur during the process
func (d *Docker) Build(containerfile string, tag string, context string) ([]byte, error) {
	args := []string{
		"build",
		"-f",
		containerfile,
		"-t",
		tag,
		context,
	}

	return runDockerCmd(args...)
}

// Exec will execute a command in the container with the provided name
// using the execOptions and the args provided. For example:
// docker exec {execOptions} {name} {args}
// Returns an error if any occur during the process
func (d *Docker) Exec(execOptions ExecOptions, name string, execArgs ...string) error {
	args := []string{
		"exec",
	}

	if execOptions.Detached {
		args = append(args, "-d")
	}

	if execOptions.Interactive {
		args = append(args, "-i")
	}

	if execOptions.Tty {
		args = append(args, "-t")
	}

	if execOptions.User != "" {
		args = append(args, "-u", execOptions.User)
	}

	if execOptions.Workdir != "" {
		args = append(args, "-w", execOptions.Workdir)
	}

	args = append(args, name)
	args = append(args, execArgs...)

	cmd := exec.Command("docker", args...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	return err
}

// ContainerList will return a list of containers
// Returns an error if any occur during the process
func (d *Docker) ContainerList() ([]Container, error) {
	containers := []Container{}
	args := []string{
		"container",
		"list",
		"--format",
		"'{{json .}}'",
	}

	out, err := runDockerCmd(args...)
	if err != nil {
		return nil, fmt.Errorf("encountered an error using `docker` to get list of containers: %w", err)
	}

	parsed := &dockerContainerList{}

	err = json.Unmarshal(out, parsed)
	if err != nil {
		return nil, fmt.Errorf("encountered an error parsing JSON from `docker container list` output: %w", err)
	}

	for _, c := range parsed.Containers {
		containers = append(containers, Container{
			Id:      c.Id,
			Name:    c.Name,
			Created: c.Created,
			Command: c.Command,
			Image:   c.Image,
			Ports:   c.Ports,
			Status:  c.Status,
			State:   c.State,
		})
	}

	return containers, nil
}

// ImageList will return a list of images
// Returns an error if any occur during the process
func (d *Docker) ImageList() ([]Image, error) {
	images := []Image{}
	args := []string{
		"image",
		"list",
		"--format",
		"'{{json .}}'",
	}

	out, err := runDockerCmd(args...)
	if err != nil {
		return nil, fmt.Errorf("encountered an error using `docker` to get list of images: %w", err)
	}

	parsed := &dockerImageList{}

	err = json.Unmarshal(out, parsed)
	if err != nil {
		return nil, fmt.Errorf("encountered an error parsing JSON from `docker image list` output: %w", err)
	}

	for _, i := range parsed.Images {
		images = append(images, Image{
			Repository: i.Repository,
			Tag:        i.Tag,
			Id:         i.ID,
			Created:    i.CreatedAt,
			Size:       i.Size,
		})
	}

	return images, nil
}

// StopContainer will stop a running container.
// Returns an error if any occur during the process
func (d *Docker) StopContainer(container Container) ([]byte, error) {
	args := []string{
		"container",
		"stop",
		container.Name,
	}

	return runDockerCmd(args...)
}

// RemoveContainer will remove a container
// Returns an error if any occur during the process
func (d *Docker) RemoveContainer(container Container) ([]byte, error) {
	args := []string{
		"container",
		"rm",
		container.Name,
	}

	return runDockerCmd(args...)
}

// CreateContainer creates a container but does not run it. Equivalent to `docker create ...`
func (d *Docker) CreateContainer(container Container) ([]byte, error) {
	args := []string{
		"create",
		"-it",
		"--name",
		container.Name,
		container.Image,
		"bash",
	}

	return runDockerCmd(args...)
}

// CopyToHost copies files from the container to the host using the provided volume.
// Returns an error if any occur during the process.
func (d *Docker) CopyToHost(container Container, volume Volume) ([]byte, error) {
	copyContainer := container
	copyContainer.Name = container.Name + "-copier"
	// create a temporary container
	out, err := d.CreateContainer(copyContainer)
	if err != nil {
		return out, fmt.Errorf("encountered an error creating temporary container to copy files: %w", err)
	}

	// copy the files
	args := []string{
		"cp",
		fmt.Sprintf("%s:%s", copyContainer.Name, volume.MountPath+"/"),
		volume.HostPath,
	}

	fmt.Println("DOCKER CP ARGS:", args)

	out, err = runDockerCmd(args...)
	if err != nil {
		return out, fmt.Errorf("encountered an error copying files: %w", err)
	}

	// remove the temporary container
	out, err = d.RemoveContainer(copyContainer)
	if err != nil {
		return out, fmt.Errorf("encountered an error removing the temporary container: %w", err)
	}

	return nil, nil
}

// runDockerCmd is a helper function to run the Docker CLI tool with the specified args.
// Returns output of the command and an error if one occurred. This blocks until command is
// complete and should not be used if you need realtime output/inputs.
func runDockerCmd(args ...string) ([]byte, error) {
	return exec.Command("docker", args...).CombinedOutput()
}
