package open

import (
	"fmt"
)

const (
	BackgroundEnvKey = "OPEN_IN_BACKGROUND"
	DisabledEnvKey   = "OPEN_DISABLED" // currently just used for simple smoke tests of the pkg
)

// Open a file, directory, or URI using the OS's default  application for that object type
func Open(uri string) error {
	cmd := open(uri)
	if cmd == nil {
		return nil
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w - unable to open uri: %s", err, output)
	}
	return nil
}

// OpenWith a file, directory, or URI using the specified application.
func OpenWith(appName, uri string) error {
	cmd := openWith(uri, appName)
	if cmd == nil {
		return nil
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w - unable to open uri with %s: %s", err, appName, output)
	}
	return nil
}
