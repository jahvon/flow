package cache

import (
	"github.com/flowexec/tuikit/io"
)

func UpdateAll(logger io.Logger) error {
	wsCache := NewWorkspaceCache()
	if err := wsCache.Update(logger); err != nil {
		return err
	}

	execCache := NewExecutableCache(wsCache)
	if err := execCache.Update(logger); err != nil {
		return err
	}

	return nil
}
