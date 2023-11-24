package cache

var (
	workspaceCache  *WorkspaceCache
	executableCache *ExecutableCache
)

func UpdateAll() error {
	wsCache := NewWorkspaceCache()
	if err := wsCache.Update(); err != nil {
		return err
	}

	execCache := NewExecutableCache()
	if err := execCache.Update(); err != nil {
		return err
	}

	return nil
}
