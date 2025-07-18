use super::command_executor::CommandExecutor;
use super::core::CommandResult;

pub struct CacheCommands<E: CommandExecutor + 'static> {
    executor: std::sync::Arc<E>,
}

impl<E: CommandExecutor + 'static> CacheCommands<E> {
    pub fn new(executor: std::sync::Arc<E>) -> Self {
        Self { executor }
    }

    pub async fn clear(&self) -> CommandResult<()> {
        self.executor.execute::<()>(&["cache", "clear"]).await
    }

    pub async fn set(&self, key: &str, value: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["cache", "set", key, value]).await
    }

    pub async fn get(&self, key: &str) -> CommandResult<String> {
        self.executor.execute(&["cache", "get", key]).await
    }

/* <<<<<<<<<<<<<<  ✨ Windsurf Command ⭐ >>>>>>>>>>>>>>>> */
    /// Remove a cache entry by key
    ///
    /// # Errors
    ///
    /// If the cache entry does not exist, or if there is an error executing the command, an error is
    /// returned.
    ///
/* <<<<<<<<<<  a303a29e-2224-47c6-9555-81d8fd250692  >>>>>>>>>>> */
    pub async fn remove(&self, key: &str) -> CommandResult<()> {
        self.executor.execute::<()>(&["cache", "remove", key]).await
    }

    pub async fn list(&self) -> CommandResult<Vec<String>> {
        self.executor.execute(&["cache", "list", "--output", "json"]).await
    }
}