//go:generate bash get_build_info.sh

package version

import (
	// using blank import for embed as it is only used inside comments.
	_ "embed"
	"fmt"
	"runtime"
	"strings"
)

var (
	// GitCommit returns the git commit that was compiled.
	//go:embed commit.txt
	gitCommit string

	// Version returns the main version number that is being exec at the moment.
	//go:embed version.txt
	version string

	// BuildDate returns the date the binary was built
	//go:embed build_date.txt
	buildDate string
)

// GoVersion returns the version of the go runtime used to compile the binary.
var goVersion = runtime.Version()

// OsArch returns the os and arch used to build the binary.
var osArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

// generateOutput return the output of the version command.
func generateOutput() string {
	return fmt.Sprintf(`

Version: %s
Git Commit: %s
Build date: %s
Go version: %s
OS / Arch : %s
`, strings.TrimSpace(version), strings.TrimSpace(gitCommit), strings.TrimSpace(buildDate), goVersion, osArch)
}

func String() string {
	return generateOutput()
}
