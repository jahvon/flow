package open

import (
	"fmt"

	open "github.com/jahvon/open-golang/open"

	"github.com/jahvon/flow/internal/io"
)

var log = io.Log().With().Str("scope", "service/open").Logger()

func Open(uri string, wait bool) error {
	log.Trace().Msgf("opening uri (%s), wait=%v", uri, wait)
	if wait {
		if err := open.Run(uri); err != nil {
			return fmt.Errorf("unable to open uri - %w", err)
		}
	} else {
		if err := open.Start(uri); err != nil {
			return fmt.Errorf("unable to open uri - %w", err)
		}
	}

	return nil
}

func OpenWith(appName, uri string, wait bool) error {
	log.Trace().Msgf("opening uri (%s) with %s, wait=%v", uri, appName, wait)
	if wait {
		if err := open.RunWith(uri, appName); err != nil {
			return fmt.Errorf("unable to open uri with %s - %w", appName, err)
		}
	} else {
		if err := open.StartWith(uri, appName); err != nil {
			return fmt.Errorf("unable to open uri with %s - %w", appName, err)
		}
	}
	return nil
}
