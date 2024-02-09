package open

import (
	"github.com/jahvon/open-golang/open"
	"github.com/pkg/errors"
)

func Open(uri string, wait bool) error {
	if wait {
		if err := open.Run(uri); err != nil {
			return errors.Wrap(err, "unable to open uri")
		}
	} else {
		if err := open.Start(uri); err != nil {
			return errors.Wrap(err, "unable to open uri")
		}
	}

	return nil
}

func OpenWith(appName, uri string, wait bool) error {
	if wait {
		if err := open.RunWith(uri, appName); err != nil {
			return errors.Wrapf(err, "unable to open uri with %s", appName)
		}
	} else {
		if err := open.StartWith(uri, appName); err != nil {
			return errors.Wrapf(err, "unable to open uri with %s", appName)
		}
	}
	return nil
}
