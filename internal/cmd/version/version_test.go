package version

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var regexp = fmt.Sprintf(`^flow - .*

Git Commit: .*
Build date: [0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} .*
Go version: go[0-9]{1}.[0-9]+.*
OS / Arch : %s %s
`, runtime.GOOS, runtime.GOARCH)

func TestGenerateOutput(t *testing.T) {
	assert.Regexp(t, regexp, generateOutput())
}

func TestPrint(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdout = w
	Print()
	w.Close()

	out, _ := io.ReadAll(r)

	assert.Regexp(t, regexp, string(out))
}
