package containerutil

// TODO(everettraven): This is meant to be used later when multiple
// Container runtimes are supported and a discovery feature is implemented
// that will return a ContainerUtil interface object corresponding to the
// proper container runtime.

// ContainerUtil is meant to generalize interactions between
// different container tools such as docker, podman, containerd, etc.
type ContainerUtil interface {
	// Run runs a container using the provided Container. It also mounts any volumes provided.
	// Returns an error if any occur during the process
	Run(container Container, volumes []Volume, runArgs ...string) ([]byte, error)

	// Build builds an image from the provided containerfile and tags it
	// with the provided tag. Returns an error if any occur during the process
	Build(containerfile string, tag string, context string) ([]byte, error)

	// Exec will execute a command in the container with the provided name
	// using the execOptions and the args provided. For example:
	// docker exec {execOptions} {name} {args}
	// Returns an error if any occur during the process
	Exec(execOptions ExecOptions, name string, execArgs ...string) error

	// ContainerList will return a list of containers
	// Returns an error if any occur during the process
	ContainerList() ([]Container, error)

	// ImageList will return a list of images
	// Returns an error if any occur during the process
	ImageList() ([]Image, error)

	// StopContainer will stop a running container.
	// Returns an error if any occur during the process
	StopContainer(container Container) ([]byte, error)

	// RemoveContainer will remove a container
	// Returns an error if any occur during the process
	RemoveContainer(container Container) ([]byte, error)

	// CopyToHost will copy files from within a container to
	// the host directory. It uses a Volume definition to determine
	// which directories to use for copy operations.
	// Returns an error if any occur during the process
	CopyToHost(container Container, volume Volume) ([]byte, error)
}

// Volume represents a Volume
type Volume struct {
	// The path on the host
	HostPath string `json:"host_path" yaml:"host_path"`
	// The path in the container
	MountPath string `json:"mount_path" yaml:"mount_path"`
}

// ExecOptions represent options that can be
// used to configure an Exec function call
type ExecOptions struct {
	Detached    bool
	Interactive bool
	Tty         bool
	User        string
	Workdir     string
}

// Container represents a container
type Container struct {
	// The container id
	Id string
	// The name of the container
	Name string
	// The image used for the container
	Image string
	// The command that was run in the container
	Command string
	// When the container was created
	Created string
	// The status of the container
	Status string
	// The state of the container
	State string
	// Ports exposed on the container
	Ports string
	// Network the container should use
	Network string
}

// Image represents an Image
type Image struct {
	// The repository/image name
	Repository string
	// The image tag
	Tag string
	// The image id
	Id string
	// When the image was created
	Created string
	// Size of the image
	Size string
}

// NewContainerUtil is used to get an implementation of ContainerUtil
// TODO(everettraven): Add discovery when multiple container runtimes
// are supported. May be nice to also add a configuration setting where
// a user can specify which runtime to be used (in the event they have
// multiple runtimes installed)
func NewContainerUtil() ContainerUtil {
	// For now since docker is the only supported runtime
	// default to returning the docker utility
	return NewDockerUtil()
}
