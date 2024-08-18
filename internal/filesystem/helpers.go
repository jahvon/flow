package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const dataDirName = "flow"

func CopyFile(src, dst string) error {
	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return errors.Wrap(err, "unable to open source file")
	}
	defer in.Close()

	data := make([]byte, 0)
	reader := bufio.NewReader(in)
	for {
		var b []byte
		b, err = reader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return errors.Wrap(err, "unable to read source file")
		}
		data = append(data, b...)
	}

	if _, err = os.Stat(dst); err == nil {
		return fmt.Errorf("file already exists: %s", dst)
	}
	if err = os.WriteFile(filepath.Clean(dst), data, 0600); err != nil {
		return errors.Wrap(err, "unable to write file")
	}

	return nil
}
