//nolint:structcheck,unused
package prm

type BackendType string

const (
	DOCKER BackendType = "docker"
)

type BackendI interface {
	GetTool(tool *Tool, prmConfig Config) error
	Validate(tool *Tool) (ToolExitCode, error)
	Exec(tool *Tool, args []string, prmConfig Config, paths DirectoryPaths) (ToolExitCode, error)
	Status() BackendStatus
}

// The BackendStatus must report whether the backend is available
// and any useful status information; in the case of the backend
// being unavailable, report the error message to the user.
type BackendStatus struct {
	IsAvailable bool
	StatusMsg   string
}

type DirectoryPaths struct {
	codeDir  string
	cacheDir string
}
