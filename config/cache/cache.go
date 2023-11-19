package cache

func UpdateAll() error {
	executableCache := NewExecutableCache()
	if err := executableCache.Update(); err != nil {
		return err
	}

	workspaceCache := NewWorkspaceCache()
	if err := workspaceCache.Update(); err != nil {
		return err
	}

	return nil
}
