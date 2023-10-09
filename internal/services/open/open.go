package open

import (
	"fmt"
	"os/exec"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log()

func Open(uri string, wait bool) error {
	log.Trace().Msgf("opening uri (%s), wait=%v", uri, wait)
	if wait {
		if err := exec.Command("open", "-W", uri).Run(); err != nil {
			return fmt.Errorf("unable to open uri - %w", err)
		}
	} else {
		if err := exec.Command("open", uri).Run(); err != nil {
			return fmt.Errorf("unable to open uri - %w", err)
		}
	}

	return nil
}

func OpenWith(appName, uri string, wait bool) error {
	log.Trace().Msgf("opening uri (%s) with %s, wait=%v", uri, appName, wait)
	if wait {
		if err := exec.Command("open", "-W", "-a", appName, uri).Run(); err != nil {
			return fmt.Errorf("unable to open uri with %s - %w", appName, err)
		}
	} else {
		if err := exec.Command("open", "-a", appName, uri).Run(); err != nil {
			return fmt.Errorf("unable to open uri with %s - %w", appName, err)
		}
	}
	return nil
}
