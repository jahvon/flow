package filesystem

import (
	cp "github.com/otiai10/copy"
)

const dataDirName = "flow"

func CopyFile(src, dst string) error {
	opts := cp.Options{
		PreserveTimes: true,
		PreserveOwner: true,
	}
	return cp.Copy(src, dst, opts)
}
