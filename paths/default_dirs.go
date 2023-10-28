package paths

const (
	defaultRootDir = "/usr/share/yet"
)

var DefaultDirs = map[string]string{
	"backups":  defaultRootDir + "/backups",
	"input":    defaultRootDir + "/input",
	"metadata": defaultRootDir + "/metadata",
	"videos":   defaultRootDir + "/videos",
}
