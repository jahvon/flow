package cache

import (
	"github.com/jahvon/tuikit/io"
)

var (
	workspaceCache  WorkspaceCache
	executableCache ExecutableCache
)

func UpdateAll(logger io.Logger) error {
	wsCache := NewWorkspaceCache()
	if err := wsCache.Update(logger); err != nil {
		return err
	}

	execCache := NewExecutableCache()
	if err := execCache.Update(logger); err != nil {
		return err
	}

	return nil
}
